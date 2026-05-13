package plc

import (
	"context"

	v1 "hotgo/api/api/plc"
	"hotgo/internal/dao"
	"hotgo/internal/service"
)

// Mines 获取当前登录用户绑定的矿场及设备列表。
func (c *ControllerV1) Mines(ctx context.Context, _ *v1.MinesReq) (res *v1.MinesRes, err error) {
	userId, err := currentFrontUserId(ctx)
	if err != nil {
		return nil, err
	}
	mineIds, err := serviceFrontUserMineIds(ctx, userId)
	if err != nil {
		return nil, err
	}
	res = &v1.MinesRes{List: []*v1.MineItem{}}
	if len(mineIds) == 0 {
		return res, nil
	}

	var mines []struct {
		Id       int    `orm:"id"`
		Name     string `orm:"name"`
		Location string `orm:"location"`
		Remark   string `orm:"remark"`
		Status   int    `orm:"status"`
	}
	err = dao.PlcMine.Ctx(ctx).
		Fields(dao.PlcMine.Columns().Id, dao.PlcMine.Columns().Name, dao.PlcMine.Columns().Location, dao.PlcMine.Columns().Remark, dao.PlcMine.Columns().Status).
		WhereIn(dao.PlcMine.Columns().Id, mineIds).
		Where(dao.PlcMine.Columns().Status, 1).
		WhereNull(dao.PlcMine.Columns().DeletedAt).
		OrderAsc(dao.PlcMine.Columns().Id).
		Scan(&mines)
	if err != nil {
		return nil, err
	}

	devicesByMine, err := userDevicesByMine(ctx, mineIds)
	if err != nil {
		return nil, err
	}
	for _, mine := range mines {
		res.List = append(res.List, &v1.MineItem{
			Id:       mine.Id,
			Name:     mine.Name,
			Location: mine.Location,
			Remark:   mine.Remark,
			Status:   mine.Status,
			Devices:  devicesByMine[mine.Id],
		})
	}
	return res, nil
}

func serviceFrontUserMineIds(ctx context.Context, userId int64) ([]int, error) {
	return service.FrontUser().MineIds(ctx, userId)
}

func userDevicesByMine(ctx context.Context, mineIds []int) (map[int][]*v1.DeviceItem, error) {
	var devices []struct {
		Id     int    `orm:"id"`
		MineId int    `orm:"mine_id"`
		Name   string `orm:"name"`
		Host   string `orm:"host"`
		Remark string `orm:"remark"`
		Status int    `orm:"status"`
	}
	err := dao.PlcDevice.Ctx(ctx).
		Fields(dao.PlcDevice.Columns().Id, dao.PlcDevice.Columns().MineId, dao.PlcDevice.Columns().Name, dao.PlcDevice.Columns().Host, dao.PlcDevice.Columns().Remark, dao.PlcDevice.Columns().Status).
		WhereIn(dao.PlcDevice.Columns().MineId, mineIds).
		Where(dao.PlcDevice.Columns().Status, 1).
		WhereNull(dao.PlcDevice.Columns().DeletedAt).
		OrderAsc(dao.PlcDevice.Columns().Id).
		Scan(&devices)
	if err != nil {
		return nil, err
	}
	out := make(map[int][]*v1.DeviceItem, len(mineIds))
	for _, d := range devices {
		out[d.MineId] = append(out[d.MineId], &v1.DeviceItem{Id: d.Id, MineId: d.MineId, Name: d.Name, Host: d.Host, Remark: d.Remark, Status: d.Status})
	}
	return out, nil
}
