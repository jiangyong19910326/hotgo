package plc

import (
	"context"

	v1 "hotgo/api/api/plc"
	"hotgo/internal/dao"
)

// Devices 获取当前启用中的设备列表
func (c *ControllerV1) Devices(ctx context.Context, _ *v1.DevicesReq) (res *v1.DevicesRes, err error) {
	userId, err := currentFrontUserId(ctx)
	if err != nil {
		return nil, err
	}
	mineIds, err := serviceFrontUserMineIds(ctx, userId)
	if err != nil {
		return nil, err
	}
	res = &v1.DevicesRes{List: []*v1.DeviceItem{}}
	if len(mineIds) == 0 {
		return res, nil
	}
	var devices []struct {
		Id     int    `orm:"id"`
		MineId int    `orm:"mine_id"`
		Name   string `orm:"name"`
		Host   string `orm:"host"`
		Remark string `orm:"remark"`
		Status int    `orm:"status"`
	}
	err = dao.PlcDevice.Ctx(ctx).
		Fields(dao.PlcDevice.Columns().Id, dao.PlcDevice.Columns().MineId, dao.PlcDevice.Columns().Name, dao.PlcDevice.Columns().Host, dao.PlcDevice.Columns().Remark, dao.PlcDevice.Columns().Status).
		WhereIn(dao.PlcDevice.Columns().MineId, mineIds).
		Where(dao.PlcDevice.Columns().Status, 1).
		WhereNull(dao.PlcDevice.Columns().DeletedAt).
		OrderAsc(dao.PlcDevice.Columns().Id).
		Scan(&devices)
	if err != nil {
		return nil, err
	}

	res.List = make([]*v1.DeviceItem, 0, len(devices))
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
