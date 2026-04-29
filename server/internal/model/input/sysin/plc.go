// Package sysin PLC 相关输入/输出模型
package sysin

import (
	"hotgo/internal/model/input/form"

	"github.com/gogf/gf/v2/os/gtime"
)

// ─────────────────────────────────────────────────────────────
// 矿场
// ─────────────────────────────────────────────────────────────

type PlcMineListInp struct {
	form.PageReq
	Name   string `json:"name"   dc:"矿场名称"`
	Status int    `json:"status" dc:"状态"`
}

type PlcMineListModel struct {
	Id        int         `json:"id"`
	Name      string      `json:"name"`
	Location  string      `json:"location"`
	Remark    string      `json:"remark"`
	Status    int         `json:"status"`
	CreatedAt *gtime.Time `json:"createdAt"`
}

type PlcMineViewInp struct {
	Id int `json:"id" v:"required#矿场ID不能为空"`
}

type PlcMineViewModel struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Remark   string `json:"remark"`
	Status   int    `json:"status"`
}

type PlcMineEditInp struct {
	Id       int    `json:"id"`
	Name     string `json:"name"     v:"required#矿场名称不能为空" dc:"矿场名称"`
	Location string `json:"location"                           dc:"地理位置"`
	Remark   string `json:"remark"                             dc:"备注"`
	Status   int    `json:"status"                             dc:"状态"`
}

type PlcMineDeleteInp struct {
	Id interface{} `json:"id" v:"required#ID不能为空"`
}

type PlcMineStatusInp struct {
	Id     int `json:"id"     v:"required#ID不能为空"`
	Status int `json:"status" v:"required#状态不能为空"`
}

// PlcMineOption 矿场下拉选项
type PlcMineOption struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// ─────────────────────────────────────────────────────────────
// PLC 设备
// ─────────────────────────────────────────────────────────────

type PlcDeviceListInp struct {
	form.PageReq
	MineId int    `json:"mineId" dc:"矿场ID"`
	Name   string `json:"name"   dc:"设备名称"`
	Status int    `json:"status" dc:"状态"`
}

type PlcDeviceListModel struct {
	Id        int         `json:"id"`
	MineId    int         `json:"mineId"`
	MineName  string      `json:"mineName"`
	Name      string      `json:"name"`
	Host      string      `json:"host"`
	Remark    string      `json:"remark"`
	Status    int         `json:"status"`
	CreatedAt *gtime.Time `json:"createdAt"`
}

type PlcDeviceViewInp struct {
	Id int `json:"id" v:"required#设备ID不能为空"`
}

type PlcDeviceViewModel struct {
	Id       int    `json:"id"`
	MineId   int    `json:"mineId"`
	MineName string `json:"mineName"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Remark   string `json:"remark"`
	Status   int    `json:"status"`
}

type PlcDeviceEditInp struct {
	Id     int    `json:"id"`
	MineId int    `json:"mineId"                                  dc:"所属矿场ID"`
	Name   string `json:"name"   v:"required#设备名称不能为空"        dc:"设备名称"`
	Host   string `json:"host"   v:"required#DTU设备编号不能为空"     dc:"DTU 设备编号 (MQTT topic)"`
	Remark string `json:"remark"                                  dc:"备注"`
	Status int    `json:"status"                                  dc:"状态"`
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
	Id        int      `json:"id"`
	DeviceId  int      `json:"deviceId"`
	Name      string   `json:"name"`
	Field     string   `json:"field"`
	DataType  string   `json:"dataType"`
	Scale     float64  `json:"scale"`
	OffsetVal float64  `json:"offsetVal"`
	Unit      string   `json:"unit"`
	AlarmMin  *float64 `json:"alarmMin"`
	AlarmMax  *float64 `json:"alarmMax"`
	Remark    string   `json:"remark"`
	Sort      int      `json:"sort"`
	Status    int      `json:"status"`
}

type PlcPointViewInp struct {
	Id int `json:"id" v:"required#ID不能为空"`
}

type PlcPointViewModel = PlcPointListModel

type PlcPointEditInp struct {
	Id        int      `json:"id"`
	DeviceId  int      `json:"deviceId"  v:"required#设备ID不能为空"      dc:"设备ID"`
	Name      string   `json:"name"      v:"required#名称不能为空"        dc:"点位名称"`
	Field     string   `json:"field"     v:"required#MQTT字段名不能为空"  dc:"MQTT 字段名 (UPPER_SNAKE)"`
	DataType  string   `json:"dataType"  v:"required#数据类型不能为空"     dc:"数据类型 Bool/Real"`
	Scale     float64  `json:"scale"                                  dc:"换算系数"`
	OffsetVal float64  `json:"offsetVal"                              dc:"换算偏移"`
	Unit      string   `json:"unit"                                   dc:"单位"`
	AlarmMin  *float64 `json:"alarmMin"                               dc:"报警下限"`
	AlarmMax  *float64 `json:"alarmMax"                               dc:"报警上限"`
	Remark    string   `json:"remark"                                 dc:"备注"`
	Sort      int      `json:"sort"                                   dc:"排序"`
	Status    int      `json:"status"                                 dc:"状态"`
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
// 应用密钥
// ─────────────────────────────────────────────────────────────

type PlcAppListInp struct {
	form.PageReq
	AppId  string `json:"appId"  dc:"AppID"`
	Name   string `json:"name"   dc:"应用名称"`
	Status int    `json:"status" dc:"状态"`
}

type PlcAppListModel struct {
	Id        int         `json:"id"`
	AppId     string      `json:"appId"`
	AppSecret string      `json:"appSecret"`
	Name      string      `json:"name"`
	Remark    string      `json:"remark"`
	Status    int         `json:"status"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

type PlcAppViewInp struct {
	Id int `json:"id" v:"required#ID不能为空"`
}

type PlcAppViewModel = PlcAppListModel

type PlcAppEditInp struct {
	Id        int    `json:"id"`
	AppId     string `json:"appId"     v:"required#AppID不能为空"     dc:"AppID"`
	AppSecret string `json:"appSecret" v:"required#AppSecret不能为空" dc:"AppSecret"`
	Name      string `json:"name"      v:"required#应用名称不能为空"   dc:"应用名称"`
	Remark    string `json:"remark"                               dc:"备注"`
	Status    int    `json:"status"                               dc:"状态"`
}

type PlcAppDeleteInp struct {
	Id interface{} `json:"id" v:"required#ID不能为空"`
}

type PlcAppStatusInp struct {
	Id     int `json:"id"     v:"required#ID不能为空"`
	Status int `json:"status" v:"required#状态不能为空"`
}

// ─────────────────────────────────────────────────────────────
// 看板汇总: 设备 + 数据点 + 实时值 一次返回
// ─────────────────────────────────────────────────────────────

type PlcOverviewInp struct {
	DeviceId int `json:"deviceId" v:"required#设备ID不能为空"`
}

type PlcOverviewPoint struct {
	PointId   int      `json:"pointId"`
	Field     string   `json:"field"`
	Name      string   `json:"name"`
	DataType  string   `json:"dataType"`
	Unit      string   `json:"unit"`
	Scale     float64  `json:"scale"`
	OffsetVal float64  `json:"offsetVal"`
	AlarmMin  *float64 `json:"alarmMin"`
	AlarmMax  *float64 `json:"alarmMax"`
	Sort      int      `json:"sort"`
	EngValue  *float64 `json:"engValue"`            // 实时工程值, nil 表示尚无数据
	AlarmType int      `json:"alarmType"`           // 0=正常 1=超上限 2=超下限
}

type PlcOverviewModel struct {
	Device *PlcDeviceViewModel `json:"device"`
	Points []*PlcOverviewPoint `json:"points"`
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
