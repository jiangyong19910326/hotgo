package plc

import (
	"context"

	v1 "hotgo/api/open/plc"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

func (c *ControllerV1) Realtime(ctx context.Context, req *v1.RealtimeReq) (res *v1.RealtimeRes, err error) {
	data, err := service.PlcRealtime().Get(ctx, &sysin.PlcRealtimeInp{DeviceId: req.DeviceId})
	if err != nil {
		return nil, err
	}

	res = &v1.RealtimeRes{
		DeviceId: data.DeviceId,
		Points:   make([]*v1.RealtimePointItem, 0, len(data.Points)),
	}
	for _, p := range data.Points {
		item := &v1.RealtimePointItem{
			PointId:   p.PointId,
			Field:     p.Field,
			Name:      p.Name,
			EngValue:  p.EngValue,
			Unit:      p.Unit,
			AlarmType: p.AlarmType,
		}
		if p.AlarmMax != nil {
			v := *p.AlarmMax
			item.AlarmMax = &v
		}
		if p.AlarmMin != nil {
			v := *p.AlarmMin
			item.AlarmMin = &v
		}
		res.Points = append(res.Points, item)
	}
	return
}
