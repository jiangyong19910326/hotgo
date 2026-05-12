package plc

import (
	"context"

	v1 "hotgo/api/api/plc"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

func (c *ControllerV1) ChartTemperature(ctx context.Context, req *v1.ChartTemperatureReq) (res *v1.ChartTemperatureRes, err error) {
	data, err := service.PlcChart().Temperature(ctx, &sysin.PlcChartTemperatureInp{
		DeviceId:     req.DeviceId,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		FieldTemp:    req.FieldTemp,
		FieldOilBack: req.FieldOilBack,
		FieldOilTank: req.FieldOilTank,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.ChartTemperatureRes{ChartModel: toV1Chart(data)}
	return
}
