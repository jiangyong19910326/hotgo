package plc

import (
	"context"
	"strings"

	v1 "hotgo/api/api/plc"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

// alarmFieldPrefix 报警点位字段前缀，所有以此开头的数据点统一归类到报警信息。
const alarmFieldPrefix = "AL_"

func (c *ControllerV1) Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error) {
	if err = ensureDeviceAccess(ctx, req.DeviceId); err != nil {
		return nil, err
	}
	data, err := service.PlcRealtime().Overview(ctx, &sysin.PlcOverviewInp{DeviceId: req.DeviceId})
	if err != nil {
		return nil, err
	}
	if data == nil || data.Device == nil {
		return &v1.OverviewRes{Points: []*v1.OverviewPoint{}}, nil
	}

	res = &v1.OverviewRes{
		Device: &v1.OverviewDevice{
			Id:       data.Device.Id,
			MineId:   data.Device.MineId,
			MineName: data.Device.MineName,
			Name:     data.Device.Name,
			Host:     data.Device.Host,
			Remark:   data.Device.Remark,
			Status:   data.Device.Status,
		},
		Points: make([]*v1.OverviewPoint, 0, len(data.Points)),
	}
	if req.WithAlarms {
		res.Alarms = make([]*v1.OverviewPoint, 0)
	}
	for _, p := range data.Points {
		item := &v1.OverviewPoint{
			PointId:     p.PointId,
			Field:       p.Field,
			Name:        p.Name,
			DataType:    p.DataType,
			Unit:        p.Unit,
			Scale:       p.Scale,
			OffsetVal:   p.OffsetVal,
			AlarmMin:    p.AlarmMin,
			AlarmMax:    p.AlarmMax,
			Sort:        p.Sort,
			EngValue:    p.EngValue,
			AlarmType:   p.AlarmType,
			CollectedAt: p.CollectedAt,
		}
		isAlarm := strings.HasPrefix(strings.ToUpper(p.Field), alarmFieldPrefix)
		fillBoolState(item, isAlarm)
		if isAlarm {
			// 实时监控页面不再返回 AL_ 点位；仅当显式 withAlarms=true 才挂到 Alarms。
			if req.WithAlarms {
				res.Alarms = append(res.Alarms, item)
			}
			continue
		}
		res.Points = append(res.Points, item)
	}
	return
}

// fillBoolState 对 Bool 点位填充语义化字段，便于前端展示运行/停止/报警/正常���
func fillBoolState(item *v1.OverviewPoint, isAlarm bool) {
	if !strings.EqualFold(item.DataType, "Bool") || item.EngValue == nil {
		return
	}
	active := *item.EngValue != 0
	item.Active = &active
	switch {
	case isAlarm && active:
		item.StateText = "报警"
	case isAlarm && !active:
		item.StateText = "正常"
	case !isAlarm && active:
		item.StateText = "运行"
	default:
		item.StateText = "停止"
	}
}
