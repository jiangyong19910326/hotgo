// Package sysin PLC 相关输入/输出模型
package sysin

import (
	"hotgo/internal/model/input/form"

	"github.com/gogf/gf/v2/os/gtime"
)

// ─────────────────────────────────────────────────────────────
// PLC 设备
// ─────────────────────────────────────────────────────────────

type PlcDeviceListInp struct {
	form.PageReq
	Name   string `json:"name"   dc:"设备名称"`
	Status int    `json:"status" dc:"状态"`
}

type PlcDeviceListModel struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	Host       string      `json:"host"`
	Port       int         `json:"port"`
	Rack       int         `json:"rack"`
	Slot       int         `json:"slot"`
	IntervalMs int         `json:"intervalMs"`
	Remark     string      `json:"remark"`
	Status     int         `json:"status"`
	CreatedAt  *gtime.Time `json:"createdAt"`
}

type PlcDeviceViewInp struct {
	Id int `json:"id" v:"required#设备ID不能为空"`
}

type PlcDeviceViewModel struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Rack       int    `json:"rack"`
	Slot       int    `json:"slot"`
	IntervalMs int    `json:"intervalMs"`
	Remark     string `json:"remark"`
	Status     int    `json:"status"`
}

type PlcDeviceEditInp struct {
	Id         int    `json:"id"`
	Name       string `json:"name"       v:"required#设备名称不能为空"    dc:"设备名称"`
	Host       string `json:"host"       v:"required#IP地址不能为空"     dc:"IP 地址"`
	Port       int    `json:"port"       v:"required#端口不能为空"       dc:"TCP 端口"`
	Rack       int    `json:"rack"                                   dc:"机架号"`
	Slot       int    `json:"slot"                                   dc:"槽号"`
	IntervalMs int    `json:"intervalMs" v:"required#采集间隔不能为空"     dc:"采集间隔(ms)"`
	Remark     string `json:"remark"                                 dc:"备注"`
	Status     int    `json:"status"                                 dc:"状态"`
}

type PlcDeviceDeleteInp struct {
	Id interface{} `json:"id" v:"required#ID不能为空"`
}

type PlcDeviceStatusInp struct {
	Id     int `json:"id"     v:"required#ID不能为空"`
	Status int `json:"status" v:"required#状态不能为空"`
}

// ─────────────────────────────────────────────────────────────
// 数据点
// ─────────────────────────────────────────────────────────────

type PlcPointListInp struct {
	form.PageReq
	DeviceId int    `json:"deviceId" dc:"设备ID"`
	Name     string `json:"name"     dc:"点位名称"`
	Status   int    `json:"status"   dc:"状态"`
}

type PlcPointListModel struct {
	Id         int      `json:"id"`
	DeviceId   int      `json:"deviceId"`
	Name       string   `json:"name"`
	Field      string   `json:"field"`
	Area       string   `json:"area"`
	DbNumber   int      `json:"dbNumber"`
	ByteOffset int      `json:"byteOffset"`
	BitOffset  int      `json:"bitOffset"`
	DataType   string   `json:"dataType"`
	Scale      float64  `json:"scale"`
	OffsetVal  float64  `json:"offsetVal"`
	Unit       string   `json:"unit"`
	AlarmMin   *float64 `json:"alarmMin"`
	AlarmMax   *float64 `json:"alarmMax"`
	Remark     string   `json:"remark"`
	Sort       int      `json:"sort"`
	Status     int      `json:"status"`
}

type PlcPointViewInp struct {
	Id int `json:"id" v:"required#ID不能为空"`
}

type PlcPointViewModel = PlcPointListModel

type PlcPointEditInp struct {
	Id         int      `json:"id"`
	DeviceId   int      `json:"deviceId"   v:"required#设备ID不能为空"   dc:"设备ID"`
	Name       string   `json:"name"       v:"required#名称不能为空"     dc:"点位名称"`
	Field      string   `json:"field"      v:"required#字段不能为空"     dc:"字段标识"`
	Area       string   `json:"area"       v:"required#存储区不能为空"    dc:"存储区"`
	DbNumber   int      `json:"dbNumber"                             dc:"DB块号"`
	ByteOffset int      `json:"byteOffset"                           dc:"字节偏移"`
	BitOffset  int      `json:"bitOffset"                            dc:"位偏移"`
	DataType   string   `json:"dataType"   v:"required#数据类型不能为空"   dc:"数据类型"`
	Scale      float64  `json:"scale"                                dc:"换算系数"`
	OffsetVal  float64  `json:"offsetVal"                            dc:"换算偏移"`
	Unit       string   `json:"unit"                                 dc:"单位"`
	AlarmMin   *float64 `json:"alarmMin"                             dc:"报警下限"`
	AlarmMax   *float64 `json:"alarmMax"                             dc:"报警上限"`
	Remark     string   `json:"remark"                               dc:"备注"`
	Sort       int      `json:"sort"                                 dc:"排序"`
	Status     int      `json:"status"                               dc:"状态"`
}

type PlcPointDeleteInp struct {
	Id interface{} `json:"id" v:"required#ID不能为空"`
}

type PlcPointStatusInp struct {
	Id     int `json:"id"     v:"required#ID不能为空"`
	Status int `json:"status" v:"required#状态不能为空"`
}

// ─────────────────────────────────────────────────────────────
// 实时数据
// ─────────────────────────────────────────────────────────────

type PlcRealtimeInp struct {
	DeviceId int `json:"deviceId" v:"required#设备ID不能为空"`
}

type PlcRealtimeItem struct {
	PointId   int      `json:"pointId"`
	Field     string   `json:"field"`
	Name      string   `json:"name"`
	EngValue  float64  `json:"engValue"`
	Unit      string   `json:"unit"`
	AlarmType int      `json:"alarmType"` // 0正常 1超上限 2超下限
	AlarmMax  *float64 `json:"alarmMax"`
	AlarmMin  *float64 `json:"alarmMin"`
}

type PlcRealtimeModel struct {
	DeviceId int                `json:"deviceId"`
	Points   []*PlcRealtimeItem `json:"points"`
}

// ─────────────────────────────────────────────────────────────
// 历史记录
// ─────────────────────────────────────────────────────────────

type PlcHistoryInp struct {
	form.PageReq
	PointId   int    `json:"pointId"   v:"required#数据点ID不能为空"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime"   dc:"结束时间"`
}

type PlcHistoryModel struct {
	CollectedAt *gtime.Time `json:"collectedAt"`
	EngValue    *float64    `json:"engValue"`
	RawValue    string      `json:"rawValue"`
}

// ─────────────────────────────────────────────────────────────
// 报警
// ─────────────────────────────────────────────────────────────

type PlcAlarmListInp struct {
	form.PageReq
	DeviceId   int    `json:"deviceId"   dc:"设备ID"`
	IsResolved int    `json:"isResolved" dc:"是否处理：1是 2否"`
	StartTime  string `json:"startTime"  dc:"开始时间"`
	EndTime    string `json:"endTime"    dc:"结束时间"`
}

type PlcAlarmListModel struct {
	Id          int64       `json:"id"`
	DeviceId    int         `json:"deviceId"`
	PointId     int         `json:"pointId"`
	PointName   string      `json:"pointName"`
	EngValue    float64     `json:"engValue"`
	AlarmType   int         `json:"alarmType"`
	AlarmMin    *float64    `json:"alarmMin"`
	AlarmMax    *float64    `json:"alarmMax"`
	Unit        string      `json:"unit"`
	IsResolved  int         `json:"isResolved"`
	ResolvedAt  *gtime.Time `json:"resolvedAt"`
	Remark      string      `json:"remark"`
	TriggeredAt *gtime.Time `json:"triggeredAt"`
}

type PlcAlarmResolveInp struct {
	Id     int64  `json:"id"     v:"required#报警ID不能为空"`
	Remark string `json:"remark" dc:"处理备注"`
}
