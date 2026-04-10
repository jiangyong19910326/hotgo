// Package sys PLC 实时数据（基于 Redis 缓存）
package sys

import (
	"context"
	"fmt"

	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sPlcRealtime struct{}

func NewPlcRealtime() *sPlcRealtime { return &sPlcRealtime{} }

func init() {
	service.RegisterPlcRealtime(NewPlcRealtime())
}

func cacheKey(deviceId int) string {
	return fmt.Sprintf("plc:realtime:%d", deviceId)
}

// Set 写入实时数据到 Redis（Hash 结构，key=field）
func (s *sPlcRealtime) Set(ctx context.Context, deviceId int, items []*sysin.PlcRealtimeItem) error {
	if len(items) == 0 {
		return nil
	}
	key := cacheKey(deviceId)
	rds, err := g.Redis().Conn(ctx)
	if err != nil {
		return err
	}
	defer rds.Close(ctx)

	for _, item := range items {
		data := g.Map{
			"pointId":   item.PointId,
			"field":     item.Field,
			"name":      item.Name,
			"engValue":  item.EngValue,
			"unit":      item.Unit,
			"alarmType": item.AlarmType,
		}
		if item.AlarmMax != nil {
			data["alarmMax"] = *item.AlarmMax
		}
		if item.AlarmMin != nil {
			data["alarmMin"] = *item.AlarmMin
		}
		jsonStr := gjson.MustEncodeString(data)
		if _, err = rds.Do(ctx, "HSET", key, item.Field, jsonStr); err != nil {
			return err
		}
	}
	// 设置过期时间：采集间隔 * 10，防止设备下线后数据无限留存
	_, _ = rds.Do(ctx, "EXPIRE", key, 600)
	return nil
}

// Get 从 Redis 读取实时数据
func (s *sPlcRealtime) Get(ctx context.Context, in *sysin.PlcRealtimeInp) (res *sysin.PlcRealtimeModel, err error) {
	key := cacheKey(in.DeviceId)
	rds, err := g.Redis().Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer rds.Close(ctx)

	val, err := rds.Do(ctx, "HGETALL", key)
	if err != nil {
		return nil, err
	}

	res = &sysin.PlcRealtimeModel{DeviceId: in.DeviceId}
	rawMap := val.MapStrStr()
	for field, jsonStr := range rawMap {
		var item sysin.PlcRealtimeItem
		if convErr := gconv.Scan(jsonStr, &item); convErr != nil {
			continue
		}
		item.Field = field
		res.Points = append(res.Points, &item)
	}
	return
}
