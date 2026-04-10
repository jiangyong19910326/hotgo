// Package crons PLC 数据采集定时任务
package crons

import (
	"context"

	"hotgo/internal/library/cron"
	"hotgo/internal/library/plc"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

func init() {
	cron.Register(PlcCollect)
}

// PlcCollect PLC 数据采集
var PlcCollect = &cPlcCollect{name: "plc_collect"}

type cPlcCollect struct {
	name string
}

func (c *cPlcCollect) GetName() string { return c.name }

// Execute 每次执行遍历所有启用设备，逐设备读取所有数据点
// 建议在系统管理 → 定时任务 中配置 Pattern 为 "* * * * * *"（每秒），
// 或按设备采集间隔要求调整。
func (c *cPlcCollect) Execute(ctx context.Context, parser *cron.Parser) (err error) {
	devices, err := service.PlcDevice().ActiveDevices(ctx)
	if err != nil {
		parser.Logger.Warningf(ctx, "plc_collect: load devices err: %v", err)
		return nil // 不中断任务
	}

	for _, dev := range devices {
		if collectErr := c.collectDevice(ctx, dev); collectErr != nil {
			parser.Logger.Warningf(ctx, "plc_collect: device[%d] %s err: %v", dev.Id, dev.Name, collectErr)
		}
	}
	return nil
}

func (c *cPlcCollect) collectDevice(ctx context.Context, dev *entity.PlcDevice) error {
	points, err := service.PlcPoint().ActivePoints(ctx, dev.Id)
	if err != nil || len(points) == 0 {
		return err
	}

	// 构建读取请求
	reqs := make([]plc.ReadRequest, 0, len(points))
	for _, p := range points {
		req := plc.ReadRequest{
			PointID:    p.Id,
			Field:      p.Field,
			Area:       plc.Area(p.Area),
			DBNumber:   p.DbNumber,
			ByteOffset: p.ByteOffset,
			BitOffset:  p.BitOffset,
			DataType:   plc.DataType(p.DataType),
			Scale:      p.Scale,
			OffsetVal:  p.OffsetVal,
			AlarmMin:   p.AlarmMin,
			AlarmMax:   p.AlarmMax,
			Unit:       p.Unit,
		}
		reqs = append(reqs, req)
	}

	// 获取或创建连接并批量读取
	client := plc.GetOrCreate(dev.Id, dev.Host, dev.Port, dev.Rack, dev.Slot)
	results := client.ReadAll(reqs)

	// 整理结果
	pointMap := make(map[int]*entity.PlcPoint, len(points))
	for _, p := range points {
		pointMap[p.Id] = p
	}

	var records []*entity.PlcRecord
	var realtimeItems []*sysin.PlcRealtimeItem
	var alarmItems []sysin.PlcRealtimeItem

	for _, r := range results {
		if !r.Valid {
			continue
		}
		engVal := r.EngValue
		records = append(records, &entity.PlcRecord{
			DeviceId: dev.Id,
			PointId:  r.PointID,
			Field:    r.Field,
			RawValue: r.RawValue,
			EngValue: &engVal,
		})

		pName := ""
		if p, ok := pointMap[r.PointID]; ok {
			pName = p.Name
		}
		item := &sysin.PlcRealtimeItem{
			PointId:   r.PointID,
			Field:     r.Field,
			Name:      pName,
			EngValue:  r.EngValue,
			Unit:      r.Unit,
			AlarmType: r.AlarmType,
			AlarmMax:  r.AlarmMax,
			AlarmMin:  r.AlarmMin,
		}
		realtimeItems = append(realtimeItems, item)
		if r.AlarmType != 0 {
			alarmItems = append(alarmItems, *item)
		}
	}

	// 写历史记录
	if len(records) > 0 {
		if insErr := service.PlcHistory().BatchInsert(ctx, records); insErr != nil {
			return insErr
		}
	}

	// 更新实时缓存
	if len(realtimeItems) > 0 {
		_ = service.PlcRealtime().Set(ctx, dev.Id, realtimeItems)
	}

	// 触发报警
	if len(alarmItems) > 0 {
		_ = service.PlcAlarm().TriggerIfNeeded(ctx, dev.Id, alarmItems)
	}

	return nil
}
