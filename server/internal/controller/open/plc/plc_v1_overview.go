package plc

import (
	"context"

	v1 "hotgo/api/open/plc"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

// Overview 看板汇总: 一次返回 设备 + 数据点骨架 + 实时值
func (c *ControllerV1) Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error) {
	data, err := service.PlcRealtime().Overview(ctx, &sysin.PlcOverviewInp{DeviceId: req.DeviceId})
	if err != nil {
		return nil, err
	}
	if data == nil || data.Device == nil {
		return &v1.OverviewRes{}, nil
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
	for _, p := range data.Points {
		op := &v1.OverviewPoint{
			PointId:   p.PointId,
			Field:     p.Field,
			Name:      p.Name,
			DataType:  p.DataType,
			Unit:      p.Unit,
			Scale:     p.Scale,
			OffsetVal: p.OffsetVal,
			AlarmMin:  p.AlarmMin,
			AlarmMax:  p.AlarmMax,
			Sort:      p.Sort,
			EngValue:  p.EngValue,
			AlarmType: p.AlarmType,
		}
		res.Points = append(res.Points, op)
	}
	return
}
