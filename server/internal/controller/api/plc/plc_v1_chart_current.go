package plc

import (
	"context"

	v1 "hotgo/api/api/plc"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

func (c *ControllerV1) ChartCurrent(ctx context.Context, req *v1.ChartCurrentReq) (res *v1.ChartCurrentRes, err error) {
	data, err := service.PlcChart().Current(ctx, &sysin.PlcChartCurrentInp{
		DeviceId:     req.DeviceId,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		FieldCurrent: req.FieldCurrent,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ChartCurrentRes{ChartModel: toV1Chart(data)}
	return
}
