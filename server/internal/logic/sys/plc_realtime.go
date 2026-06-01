// Package sys PLC 实时数据（基于 Redis 缓存）
package sys

import (
	"context"
	"fmt"

	"hotgo/internal/dao"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
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

	now := gtime.Now()
	for _, item := range items {
		if item.CollectedAt == nil {
			item.CollectedAt = now
		}
		data := g.Map{
			"pointId":     item.PointId,
			"field":       item.Field,
			"name":        item.Name,
			"engValue":    item.EngValue,
			"unit":        item.Unit,
			"alarmType":   item.AlarmType,
			"collectedAt": item.CollectedAt,
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
	// 不设过期: 设备掉线时仍展示最后状态, 删设备/手动 DEL 才清
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

// Overview 看板汇总: 设备信息 + 全部启用点位 (按sort排序) + 实时值合并
// Redis 实时缓存命中优先; 缓存miss的点位用 hg_plc_record 最新一条 DB fallback
func (s *sPlcRealtime) Overview(ctx context.Context, in *sysin.PlcOverviewInp) (res *sysin.PlcOverviewModel, err error) {
	devView, err := service.PlcDevice().View(ctx, &sysin.PlcDeviceViewInp{Id: in.DeviceId})
	if err != nil {
		return nil, err
	}
	if devView == nil || devView.Id == 0 {
		return nil, nil
	}

	points, err := service.PlcPoint().ActivePoints(ctx, in.DeviceId)
	if err != nil {
		return nil, err
	}

	// 拉一次实时缓存
	rt, err := s.Get(ctx, &sysin.PlcRealtimeInp{DeviceId: in.DeviceId})
	if err != nil {
		return nil, err
	}
	rtMap := make(map[string]*sysin.PlcRealtimeItem, len(rt.Points))
	for _, p := range rt.Points {
		rtMap[p.Field] = p
	}

	// 收集缓存miss的pointId, 一次批量从record表取最新
	var missIds []int
	for _, p := range points {
		if _, ok := rtMap[p.Field]; !ok {
			missIds = append(missIds, p.Id)
		}
	}
	dbFallback := make(map[int]struct {
		EngValue    float64
		CollectedAt *gtime.Time
	}) // pointId -> latest record
	if len(missIds) > 0 {
		var rows []struct {
			PointId     int         `orm:"point_id"`
			EngValue    *float64    `orm:"eng_value"`
			CollectedAt *gtime.Time `orm:"collected_at"`
		}
		// 子查询: 每个 point 的最大 collected_at 对应一条
		err = dao.PlcRecord.Ctx(ctx).
			Fields("point_id, eng_value, collected_at").
			Where("id IN (?)",
				dao.PlcRecord.Ctx(ctx).
					Fields("MAX(id) as id").
					WhereIn("point_id", missIds).
					Group("point_id"),
			).Scan(&rows)
		if err == nil {
			for _, r := range rows {
				if r.EngValue != nil {
					dbFallback[r.PointId] = struct {
						EngValue    float64
						CollectedAt *gtime.Time
					}{EngValue: *r.EngValue, CollectedAt: r.CollectedAt}
				}
			}
		}
	}

	out := &sysin.PlcOverviewModel{
		Device: devView,
		Points: make([]*sysin.PlcOverviewPoint, 0, len(points)),
	}
	for _, p := range points {
		op := &sysin.PlcOverviewPoint{
			PointId:   p.Id,
			Field:     p.Field,
			Name:      p.Name,
			DataType:  p.DataType,
			Unit:      p.Unit,
			Scale:     p.Scale,
			OffsetVal: p.OffsetVal,
			AlarmMin:  p.AlarmMin,
			AlarmMax:  p.AlarmMax,
			Sort:      p.Sort,
		}
		if r, ok := rtMap[p.Field]; ok {
			v := r.EngValue
			op.EngValue = &v
			op.AlarmType = r.AlarmType
			op.CollectedAt = r.CollectedAt
		} else if v, ok := dbFallback[p.Id]; ok {
			vv := v.EngValue
			op.EngValue = &vv
			op.CollectedAt = v.CollectedAt
			// 报警判定: 仅用阈值, 不重写
			if p.AlarmMax != nil && vv > *p.AlarmMax {
				op.AlarmType = 1
			} else if p.AlarmMin != nil && vv < *p.AlarmMin {
				op.AlarmType = 2
			}
		}
		out.Points = append(out.Points, op)
	}
	return out, nil
}
