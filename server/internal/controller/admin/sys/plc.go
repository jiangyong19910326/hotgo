// Package sys PLC 管理控制器
package sys

import (
	"context"

	"hotgo/api/admin/plc"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

var Plc = cPlc{}

type cPlc struct{}

// ─────────────────────────────────────────────────────────────
// 设备
// ─────────────────────────────────────────────────────────────

func (c *cPlc) DeviceList(ctx context.Context, req *plc.DeviceListReq) (res *plc.DeviceListRes, err error) {
	list, total, err := service.PlcDevice().List(ctx, &req.PlcDeviceListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PlcDeviceListModel{}
	}
	res = new(plc.DeviceListRes)
	res.List = list
	res.PageRes.Pack(req, total)
	return
}

func (c *cPlc) DeviceView(ctx context.Context, req *plc.DeviceViewReq) (res *plc.DeviceViewRes, err error) {
	data, err := service.PlcDevice().View(ctx, &req.PlcDeviceViewInp)
	if err != nil {
		return
	}
	res = new(plc.DeviceViewRes)
	res.PlcDeviceViewModel = data
	return
}

func (c *cPlc) DeviceEdit(ctx context.Context, req *plc.DeviceEditReq) (res *plc.DeviceEditRes, err error) {
	err = service.PlcDevice().Edit(ctx, &req.PlcDeviceEditInp)
	return
}

func (c *cPlc) DeviceDelete(ctx context.Context, req *plc.DeviceDeleteReq) (res *plc.DeviceDeleteRes, err error) {
	err = service.PlcDevice().Delete(ctx, &req.PlcDeviceDeleteInp)
	return
}

func (c *cPlc) DeviceStatus(ctx context.Context, req *plc.DeviceStatusReq) (res *plc.DeviceStatusRes, err error) {
	err = service.PlcDevice().Status(ctx, &req.PlcDeviceStatusInp)
	return
}

// ─────────────────────────────────────────────────────────────
// 数据点
// ─────────────────────────────────────────────────────────────

func (c *cPlc) PointList(ctx context.Context, req *plc.PointListReq) (res *plc.PointListRes, err error) {
	list, total, err := service.PlcPoint().List(ctx, &req.PlcPointListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PlcPointListModel{}
	}
	res = new(plc.PointListRes)
	res.List = list
	res.PageRes.Pack(req, total)
	return
}

func (c *cPlc) PointView(ctx context.Context, req *plc.PointViewReq) (res *plc.PointViewRes, err error) {
	data, err := service.PlcPoint().View(ctx, &req.PlcPointViewInp)
	if err != nil {
		return
	}
	res = new(plc.PointViewRes)
	res.PlcPointViewModel = data
	return
}

func (c *cPlc) PointEdit(ctx context.Context, req *plc.PointEditReq) (res *plc.PointEditRes, err error) {
	err = service.PlcPoint().Edit(ctx, &req.PlcPointEditInp)
	return
}

func (c *cPlc) PointDelete(ctx context.Context, req *plc.PointDeleteReq) (res *plc.PointDeleteRes, err error) {
	err = service.PlcPoint().Delete(ctx, &req.PlcPointDeleteInp)
	return
}

func (c *cPlc) PointStatus(ctx context.Context, req *plc.PointStatusReq) (res *plc.PointStatusRes, err error) {
	err = service.PlcPoint().Status(ctx, &req.PlcPointStatusInp)
	return
}

// ─────────────────────────────────────────────────────────────
// 实时 / 历史 / 报警
// ─────────────────────────────────────────────────────────────

func (c *cPlc) Realtime(ctx context.Context, req *plc.RealtimeReq) (res *plc.RealtimeRes, err error) {
	data, err := service.PlcRealtime().Get(ctx, &req.PlcRealtimeInp)
	if err != nil {
		return
	}
	res = new(plc.RealtimeRes)
	res.PlcRealtimeModel = data
	return
}

func (c *cPlc) History(ctx context.Context, req *plc.HistoryReq) (res *plc.HistoryRes, err error) {
	list, total, err := service.PlcHistory().List(ctx, &req.PlcHistoryInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PlcHistoryModel{}
	}
	res = new(plc.HistoryRes)
	res.List = list
	res.PageRes.Pack(req, total)
	return
}

func (c *cPlc) AlarmList(ctx context.Context, req *plc.AlarmListReq) (res *plc.AlarmListRes, err error) {
	list, total, err := service.PlcAlarm().List(ctx, &req.PlcAlarmListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PlcAlarmListModel{}
	}
	res = new(plc.AlarmListRes)
	res.List = list
	res.PageRes.Pack(req, total)
	return
}

func (c *cPlc) AlarmResolve(ctx context.Context, req *plc.AlarmResolveReq) (res *plc.AlarmResolveRes, err error) {
	err = service.PlcAlarm().Resolve(ctx, &req.PlcAlarmResolveInp)
	return
}
