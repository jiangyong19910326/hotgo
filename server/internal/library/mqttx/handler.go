// Package mqttx MQTT 消息处理
package mqttx

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
)

// payloadItem MQTT 单字段载荷
type payloadItem struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

// onMessage MQTT 消息回调
// 流程: 解析topic取deviceCode → 解析payload → 查/建device → 查/建points → 入库/缓存/报警
func onMessage(_ mqtt.Client, msg mqtt.Message) {
	ctx := context.Background()
	topic := msg.Topic()

	deviceCode := extractDeviceCode(topic)
	if deviceCode == "" {
		g.Log().Warningf(ctx, "mqttx: unexpected topic format: %s", topic)
		return
	}

	var payload map[string]payloadItem
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		g.Log().Warningf(ctx, "mqttx: parse payload for %s err: %v", deviceCode, err)
		return
	}
	if len(payload) == 0 {
		return
	}

	dev, err := GetOrCreateDevice(ctx, deviceCode)
	if err != nil {
		g.Log().Warningf(ctx, "mqttx: get/create device %s err: %v", deviceCode, err)
		return
	}

	fieldVals := make(map[string]float64, len(payload))
	for k, v := range payload {
		fieldVals[k] = v.Value
	}
	points, err := EnsurePoints(ctx, dev.Id, fieldVals)
	if err != nil {
		g.Log().Warningf(ctx, "mqttx: ensure points for device %d err: %v", dev.Id, err)
		return
	}
	if len(points) == 0 {
		return
	}

	records, items, alarms := buildResults(dev, points, payload)

	if len(records) > 0 {
		if insErr := service.PlcHistory().BatchInsert(ctx, records); insErr != nil {
			g.Log().Warningf(ctx, "mqttx: BatchInsert err: %v", insErr)
		}
	}
	if len(items) > 0 {
		if rtErr := service.PlcRealtime().Set(ctx, dev.Id, items); rtErr != nil {
			g.Log().Warningf(ctx, "mqttx: Realtime.Set err: %v", rtErr)
		}
	}
	if len(alarms) > 0 {
		if aErr := service.PlcAlarm().TriggerIfNeeded(ctx, dev.Id, alarms); aErr != nil {
			g.Log().Warningf(ctx, "mqttx: Alarm.Trigger err: %v", aErr)
		}
	}
}

// extractDeviceCode 从 topic /dtu/{code}/data 取 code
func extractDeviceCode(topic string) string {
	parts := strings.Split(strings.Trim(topic, "/"), "/")
	if len(parts) < 3 {
		return ""
	}
	return parts[len(parts)-2]
}

// buildResults 按 points 定义遍历 payload, 生成 record/realtime/alarm
func buildResults(dev *entity.PlcDevice, points []*entity.PlcPoint, payload map[string]payloadItem) (
	records []*entity.PlcRecord,
	items []*sysin.PlcRealtimeItem,
	alarms []sysin.PlcRealtimeItem,
) {
	for _, p := range points {
		raw, ok := payload[p.Field]
		if !ok {
			continue
		}

		engVal := raw.Value*p.Scale + p.OffsetVal
		engVal = math.Round(engVal*10000) / 10000

		alarmType := 0
		if p.AlarmMax != nil && engVal > *p.AlarmMax {
			alarmType = 1
		} else if p.AlarmMin != nil && engVal < *p.AlarmMin {
			alarmType = 2
		}

		engCopy := engVal
		records = append(records, &entity.PlcRecord{
			DeviceId: dev.Id,
			PointId:  p.Id,
			Field:    p.Field,
			RawValue: strconv.FormatFloat(raw.Value, 'f', -1, 64),
			EngValue: &engCopy,
		})

		item := &sysin.PlcRealtimeItem{
			PointId:   p.Id,
			Field:     p.Field,
			Name:      p.Name,
			EngValue:  engVal,
			Unit:      p.Unit,
			AlarmType: alarmType,
			AlarmMax:  p.AlarmMax,
			AlarmMin:  p.AlarmMin,
		}
		items = append(items, item)
		if alarmType != 0 {
			alarms = append(alarms, *item)
		}
	}
	return
}
