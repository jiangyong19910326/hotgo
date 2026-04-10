// Package plc PLC 管理接口
package plc

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

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
