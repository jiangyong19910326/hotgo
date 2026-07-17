// Package sys PLC 管理控制器
package sys

import (
	"context"
	"strings"

	"hotgo/api/admin/plc"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

var Plc = cPlc{}

type cPlc struct{}

// ─────────────────────────────────────────────────────────────
// 矿场
// ─────────────────────────────────────────────────────────────

func (c *cPlc) MineList(ctx context.Context, req *plc.MineListReq) (res *plc.MineListRes, err error) {
	list, total, err := service.PlcMine().List(ctx, &req.PlcMineListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PlcMineListModel{}
	}
	res = new(plc.MineListRes)
	res.List = list
	res.PageRes.Pack(req, total)
	return
}

func (c *cPlc) MineView(ctx context.Context, req *plc.MineViewReq) (res *plc.MineViewRes, err error) {
	data, err := service.PlcMine().View(ctx, &req.PlcMineViewInp)
	if err != nil {
		return
	}
	res = new(plc.MineViewRes)
	res.PlcMineViewModel = data
	return
}

func (c *cPlc) MineEdit(ctx context.Context, req *plc.MineEditReq) (res *plc.MineEditRes, err error) {
	err = service.PlcMine().Edit(ctx, &req.PlcMineEditInp)
	return
}

func (c *cPlc) MineDelete(ctx context.Context, req *plc.MineDeleteReq) (res *plc.MineDeleteRes, err error) {
	err = service.PlcMine().Delete(ctx, &req.PlcMineDeleteInp)
	return
}

func (c *cPlc) MineStatus(ctx context.Context, req *plc.MineStatusReq) (res *plc.MineStatusRes, err error) {
	err = service.PlcMine().Status(ctx, &req.PlcMineStatusInp)
	return
}

func (c *cPlc) MineOptions(ctx context.Context, req *plc.MineOptionsReq) (res *plc.MineOptionsRes, err error) {
	list, err := service.PlcMine().Options(ctx)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PlcMineOption{}
	}
	res = new(plc.MineOptionsRes)
	res.List = list
	return
}

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

func (c *cPlc) DeviceControl(ctx context.Context, req *plc.DeviceControlReq) (res *plc.DeviceControlRes, err error) {
	err = service.PlcDevice().Control(ctx, &req.PlcDeviceControlInp)
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

func (c *cPlc) Overview(ctx context.Context, req *plc.OverviewReq) (res *plc.OverviewRes, err error) {
	data, err := service.PlcRealtime().Overview(ctx, &req.PlcOverviewInp)
	if err != nil {
		return
	}
	res = &plc.OverviewRes{Points: []*plc.OverviewPointVO{}}
	if data == nil || data.Device == nil {
		return
	}
	res.Device = data.Device
	if req.WithAlarms {
		res.Alarms = []*plc.OverviewPointVO{}
	}
	for _, p := range data.Points {
		vo := &plc.OverviewPointVO{PlcOverviewPoint: p}
		isAlarm := strings.HasPrefix(strings.ToUpper(p.Field), "AL_")
		fillOverviewBoolState(vo, isAlarm)
		if isAlarm {
			if req.WithAlarms {
				res.Alarms = append(res.Alarms, vo)
			}
			continue
		}
		res.Points = append(res.Points, vo)
	}
	return
}

// fillOverviewBoolState 给 Bool 点位补 active / stateText，便于前端展示运行/停止/报警/正常。
func fillOverviewBoolState(vo *plc.OverviewPointVO, isAlarm bool) {
	if vo == nil || vo.PlcOverviewPoint == nil {
		return
	}
	if !strings.EqualFold(vo.DataType, "Bool") || vo.EngValue == nil {
		return
	}
	active := *vo.EngValue != 0
	vo.Active = &active
	switch {
	case isAlarm && active:
		vo.StateText = "报警"
	case isAlarm && !active:
		vo.StateText = "正常"
	case !isAlarm && active:
		vo.StateText = "运行"
	default:
		vo.StateText = "停止"
	}
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

// ─────────────────────────────────────────────────────────────
// 应用密钥
// ─────────────────────────────────────────────────────────────

func (c *cPlc) AppList(ctx context.Context, req *plc.AppListReq) (res *plc.AppListRes, err error) {
	list, total, err := service.PlcApp().List(ctx, &req.PlcAppListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PlcAppListModel{}
	}
	res = new(plc.AppListRes)
	res.List = list
	res.PageRes.Pack(req, total)
	return
}

func (c *cPlc) AppView(ctx context.Context, req *plc.AppViewReq) (res *plc.AppViewRes, err error) {
	data, err := service.PlcApp().View(ctx, &req.PlcAppViewInp)
	if err != nil {
		return
	}
	res = new(plc.AppViewRes)
	res.PlcAppViewModel = data
	return
}

func (c *cPlc) AppEdit(ctx context.Context, req *plc.AppEditReq) (res *plc.AppEditRes, err error) {
	err = service.PlcApp().Edit(ctx, &req.PlcAppEditInp)
	return
}

func (c *cPlc) AppDelete(ctx context.Context, req *plc.AppDeleteReq) (res *plc.AppDeleteRes, err error) {
	err = service.PlcApp().Delete(ctx, &req.PlcAppDeleteInp)
	return
}

func (c *cPlc) AppStatus(ctx context.Context, req *plc.AppStatusReq) (res *plc.AppStatusRes, err error) {
	err = service.PlcApp().Status(ctx, &req.PlcAppStatusInp)
	return
}

func (c *cPlc) AppGenSecret(ctx context.Context, _ *plc.AppGenSecretReq) (res *plc.AppGenSecretRes, err error) {
	res = &plc.AppGenSecretRes{AppSecret: service.PlcApp().GenSecret(ctx)}
	return
}

// ─────────────────────────────────────────────────────────────
// 图表数据（ECharts）
// ─────────────────────────────────────────────────────────────

func (c *cPlc) ChartTemperature(ctx context.Context, req *plc.ChartTemperatureReq) (res *plc.ChartTemperatureRes, err error) {
	data, err := service.PlcChart().Temperature(ctx, &req.PlcChartTemperatureInp)
	if err != nil {
		return
	}
	res = &plc.ChartTemperatureRes{PlcChartModel: data}
	return
}

func (c *cPlc) ChartCurrent(ctx context.Context, req *plc.ChartCurrentReq) (res *plc.ChartCurrentRes, err error) {
	data, err := service.PlcChart().Current(ctx, &req.PlcChartCurrentInp)
	if err != nil {
		return
	}
	res = &plc.ChartCurrentRes{PlcChartModel: data}
	return
}
