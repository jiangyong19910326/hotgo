// Package plc PLC 管理接口
package plc

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

// ─────────────────────────────────────────────────────────────
// 矿场管理
// ─────────────────────────────────────────────────────────────

type MineListReq struct {
	g.Meta `path:"/plc/mine/list" method:"get" tags:"矿场管理" summary:"获取矿场列表"`
	sysin.PlcMineListInp
}
type MineListRes struct {
	form.PageRes
	List []*sysin.PlcMineListModel `json:"list"`
}

type MineViewReq struct {
	g.Meta `path:"/plc/mine/view" method:"get" tags:"矿场管理" summary:"获取矿场详情"`
	sysin.PlcMineViewInp
}
type MineViewRes struct {
	*sysin.PlcMineViewModel
}

type MineEditReq struct {
	g.Meta `path:"/plc/mine/edit" method:"post" tags:"矿场管理" summary:"新增/修改矿场"`
	sysin.PlcMineEditInp
}
type MineEditRes struct{}

type MineDeleteReq struct {
	g.Meta `path:"/plc/mine/delete" method:"post" tags:"矿场管理" summary:"删除矿场"`
	sysin.PlcMineDeleteInp
}
type MineDeleteRes struct{}

type MineStatusReq struct {
	g.Meta `path:"/plc/mine/status" method:"post" tags:"矿场管理" summary:"更新矿场状态"`
	sysin.PlcMineStatusInp
}
type MineStatusRes struct{}

type MineOptionsReq struct {
	g.Meta `path:"/plc/mine/options" method:"get" tags:"矿场管理" summary:"获取矿场下拉选项"`
}
type MineOptionsRes struct {
	List []*sysin.PlcMineOption `json:"list"`
}

// ─────────────────────────────────────────────────────────────
// 设备管理
// ─────────────────────────────────────────────────────────────

type DeviceListReq struct {
	g.Meta `path:"/plc/device/list" method:"get" tags:"PLC设备" summary:"获取PLC设备列表"`
	sysin.PlcDeviceListInp
}
type DeviceListRes struct {
	form.PageRes
	List []*sysin.PlcDeviceListModel `json:"list"`
}

type DeviceViewReq struct {
	g.Meta `path:"/plc/device/view" method:"get" tags:"PLC设备" summary:"获取PLC设备详情"`
	sysin.PlcDeviceViewInp
}
type DeviceViewRes struct {
	*sysin.PlcDeviceViewModel
}

type DeviceEditReq struct {
	g.Meta `path:"/plc/device/edit" method:"post" tags:"PLC设备" summary:"新增/修改PLC设备"`
	sysin.PlcDeviceEditInp
}
type DeviceEditRes struct{}

type DeviceDeleteReq struct {
	g.Meta `path:"/plc/device/delete" method:"post" tags:"PLC设备" summary:"删除PLC设备"`
	sysin.PlcDeviceDeleteInp
}
type DeviceDeleteRes struct{}

type DeviceStatusReq struct {
	g.Meta `path:"/plc/device/status" method:"post" tags:"PLC设备" summary:"更新PLC设备状态"`
	sysin.PlcDeviceStatusInp
}
type DeviceStatusRes struct{}

type DeviceControlReq struct {
	g.Meta `path:"/plc/device/control" method:"post" tags:"PLC设备" summary:"发送PLC设备启停命令"`
	sysin.PlcDeviceControlInp
}
type DeviceControlRes struct {
	*sysin.PlcDeviceControlModel
}

// ─────────────────────────────────────────────────────────────
// 数据点
// ─────────────────────────────────────────────────────────────

type PointListReq struct {
	g.Meta `path:"/plc/point/list" method:"get" tags:"PLC数据点" summary:"获取数据点列表"`
	sysin.PlcPointListInp
}
type PointListRes struct {
	form.PageRes
	List []*sysin.PlcPointListModel `json:"list"`
}

type PointViewReq struct {
	g.Meta `path:"/plc/point/view" method:"get" tags:"PLC数据点" summary:"获取数据点详情"`
	sysin.PlcPointViewInp
}
type PointViewRes struct {
	*sysin.PlcPointViewModel
}

type PointEditReq struct {
	g.Meta `path:"/plc/point/edit" method:"post" tags:"PLC数据点" summary:"新增/修改数据点"`
	sysin.PlcPointEditInp
}
type PointEditRes struct{}

type PointDeleteReq struct {
	g.Meta `path:"/plc/point/delete" method:"post" tags:"PLC数据点" summary:"删除数据点"`
	sysin.PlcPointDeleteInp
}
type PointDeleteRes struct{}

type PointStatusReq struct {
	g.Meta `path:"/plc/point/status" method:"post" tags:"PLC数据点" summary:"更新数据点状态"`
	sysin.PlcPointStatusInp
}
type PointStatusRes struct{}

// ─────────────────────────────────────────────────────────────
// 实时数据
// ─────────────────────────────────────────────────────────────

type RealtimeReq struct {
	g.Meta `path:"/plc/realtime" method:"get" tags:"PLC实时" summary:"获取实时数据"`
	sysin.PlcRealtimeInp
}
type RealtimeRes struct {
	*sysin.PlcRealtimeModel
}

type OverviewReq struct {
	g.Meta `path:"/plc/overview" method:"get" tags:"PLC实时" summary:"看板汇总(设备+数据点+实时值)"`
	sysin.PlcOverviewInp
	WithAlarms bool `p:"withAlarms" dc:"是否同时返回报警点(AL_前缀)，默认 false"`
}

// OverviewPointVO 实时监控用的扩展点位（增加 active / stateText）。
type OverviewPointVO struct {
	*sysin.PlcOverviewPoint
	Active    *bool  `json:"active,omitempty"`
	StateText string `json:"stateText,omitempty"`
}

type OverviewRes struct {
	Device *sysin.PlcDeviceViewModel `json:"device"`
	Points []*OverviewPointVO        `json:"points"`           // 普通数据点，实时监控页面用
	Alarms []*OverviewPointVO        `json:"alarms,omitempty"` // 报警点（AL_前缀），仅 withAlarms=true 时返回
}

// ─────────────────────────────────────────────────────────────
// 历史记录
// ─────────────────────────────────────────────────────────────

type HistoryReq struct {
	g.Meta `path:"/plc/history" method:"get" tags:"PLC历史" summary:"获取历史记录"`
	sysin.PlcHistoryInp
}
type HistoryRes struct {
	form.PageRes
	List []*sysin.PlcHistoryModel `json:"list"`
}

// ─────────────────────────────────────────────────────────────
// 报警
// ─────────────────────────────────────────────────────────────

type AlarmListReq struct {
	g.Meta `path:"/plc/alarm/list" method:"get" tags:"PLC报警" summary:"获取报警列表"`
	sysin.PlcAlarmListInp
}
type AlarmListRes struct {
	form.PageRes
	List []*sysin.PlcAlarmListModel `json:"list"`
}

type AlarmResolveReq struct {
	g.Meta `path:"/plc/alarm/resolve" method:"post" tags:"PLC报警" summary:"处理报警"`
	sysin.PlcAlarmResolveInp
}
type AlarmResolveRes struct{}

// ─────────────────────────────────────────────────────────────
// 应用密钥
// ─────────────────────────────────────────────────────────────

type AppListReq struct {
	g.Meta `path:"/plc/app/list" method:"get" tags:"PLC应用密钥" summary:"获取应用密钥列表"`
	sysin.PlcAppListInp
}
type AppListRes struct {
	form.PageRes
	List []*sysin.PlcAppListModel `json:"list"`
}

type AppViewReq struct {
	g.Meta `path:"/plc/app/view" method:"get" tags:"PLC应用密钥" summary:"获取应用密钥详情"`
	sysin.PlcAppViewInp
}
type AppViewRes struct {
	*sysin.PlcAppViewModel
}

type AppEditReq struct {
	g.Meta `path:"/plc/app/edit" method:"post" tags:"PLC应用密钥" summary:"新增/修改应用密钥"`
	sysin.PlcAppEditInp
}
type AppEditRes struct{}

type AppDeleteReq struct {
	g.Meta `path:"/plc/app/delete" method:"post" tags:"PLC应用密钥" summary:"删除应用密钥"`
	sysin.PlcAppDeleteInp
}
type AppDeleteRes struct{}

type AppStatusReq struct {
	g.Meta `path:"/plc/app/status" method:"post" tags:"PLC应用密钥" summary:"更新应用密钥状态"`
	sysin.PlcAppStatusInp
}
type AppStatusRes struct{}

type AppGenSecretReq struct {
	g.Meta `path:"/plc/app/genSecret" method:"get" tags:"PLC应用密钥" summary:"生成随机AppSecret"`
}
type AppGenSecretRes struct {
	AppSecret string `json:"appSecret"`
}

// ─────────────────────────────────────────────────────────────
// 图表数据（ECharts，监控大屏用）
// ─────────────────────────────────────────────────────────────

type ChartTemperatureReq struct {
	g.Meta `path:"/plc/chart/temperature" method:"get" tags:"PLC图表" summary:"设备温度折线图(设备/回油/油箱)"`
	sysin.PlcChartTemperatureInp
}
type ChartTemperatureRes struct {
	*sysin.PlcChartModel
}

type ChartCurrentReq struct {
	g.Meta `path:"/plc/chart/current" method:"get" tags:"PLC图表" summary:"设备电流折线图"`
	sysin.PlcChartCurrentInp
}
type ChartCurrentRes struct {
	*sysin.PlcChartModel
}
