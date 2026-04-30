package plc

import (
	"context"

	v1 "hotgo/api/api/plc"
	"hotgo/internal/service"
)

// Devices 获取当前启用中的设备列表
func (c *ControllerV1) Devices(ctx context.Context, _ *v1.DevicesReq) (res *v1.DevicesRes, err error) {
	devices, err := service.PlcDevice().ActiveDevices(ctx)
	if err != nil {
		return nil, err
	}

	res = &v1.DevicesRes{List: make([]*v1.DeviceItem, 0, len(devices))}
	for _, d := range devices {
		res.List = append(res.List, &v1.DeviceItem{
			Id:     d.Id,
			MineId: d.MineId,
			Name:   d.Name,
			Host:   d.Host,
			Remark: d.Remark,
			Status: d.Status,
		})
	}
	return
}
