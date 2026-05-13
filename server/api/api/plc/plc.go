// Package plc 前台 PLC 接口（暂未启用签名验签）
package plc

import "github.com/gogf/gf/v2/frame/g"

// DevicesReq 当前可用设备列表
type DevicesReq struct {
	g.Meta `path:"/plc/devices" method:"get" tags:"PLC前台" summary:"获取当前启用中的设备列表"`
}

type DeviceItem struct {
	Id     int    `json:"id"`
	MineId int    `json:"mineId"`
	Name   string `json:"name"`
	Host   string `json:"host"`
	Remark string `json:"remark"`
	Status int    `json:"status"`
}

type DevicesRes struct {
	List []*DeviceItem `json:"list"`
}

// MinesReq 当前登录用户绑定的矿场及设备列表
type MinesReq struct {
	g.Meta `path:"/plc/mines" method:"get" tags:"PLC前台" summary:"获取当前用户绑定的矿场及设备列表"`
}

type MineItem struct {
	Id       int           `json:"id"`
	Name     string        `json:"name"`
	Location string        `json:"location"`
	Remark   string        `json:"remark"`
	Status   int           `json:"status"`
	Devices  []*DeviceItem `json:"devices"`
}

type MinesRes struct {
	List []*MineItem `json:"list"`
}

// OverviewReq 看板汇总: 设备 + 数据点 + 实时值
type OverviewReq struct {
	g.Meta   `path:"/plc/overview" method:"get" tags:"PLC前台" summary:"看板汇总(设备+数据点+实时值)"`
	DeviceId int `p:"deviceId" v:"required|min:1#设备ID不能为空" dc:"设备ID"`
}

type OverviewDevice struct {
	Id       int    `json:"id"`
	MineId   int    `json:"mineId"`
	MineName string `json:"mineName"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Remark   string `json:"remark"`
	Status   int    `json:"status"`
}

type OverviewPoint struct {
	PointId   int      `json:"pointId"`
	Field     string   `json:"field"`
	Name      string   `json:"name"`
	DataType  string   `json:"dataType"`
	Unit      string   `json:"unit"`
	Scale     float64  `json:"scale"`
	OffsetVal float64  `json:"offsetVal"`
	AlarmMin  *float64 `json:"alarmMin,omitempty"`
	AlarmMax  *float64 `json:"alarmMax,omitempty"`
	Sort      int      `json:"sort"`
	EngValue  *float64 `json:"engValue"`
	AlarmType int      `json:"alarmType"`
}

type OverviewRes struct {
	Device *OverviewDevice  `json:"device"`
	Points []*OverviewPoint `json:"points"`
}

// HistoryReq 单点位历史数据 (画图用)
type HistoryReq struct {
	g.Meta    `path:"/plc/history" method:"get" tags:"PLC前台" summary:"数据点历史值(画图)"`
	DeviceId  int    `p:"deviceId"  v:"required|min:1#设备ID不能为空" dc:"设备ID"`
	PointId   int    `p:"pointId"  v:"required|min:1#点位ID不能为空" dc:"点位ID"`
	StartTime string `p:"startTime" dc:"开始时间, 格式 2006-01-02 15:04:05, 为空默认 endTime-1h"`
	EndTime   string `p:"endTime"   dc:"结束时间, 为空默认现在"`
	Limit     int    `p:"limit"     d:"100" dc:"返回点数上限, 默认100, 最大5000"`
}

type HistoryPointInfo struct {
	PointId  int    `json:"pointId"`
	DeviceId int    `json:"deviceId"`
	Field    string `json:"field"`
	Name     string `json:"name"`
	Unit     string `json:"unit"`
}

type HistoryItem struct {
	Time     string   `json:"time"`     // 采集时间 2006-01-02 15:04:05
	EngValue *float64 `json:"engValue"` // 工程值
}

type HistoryRes struct {
	Point *HistoryPointInfo `json:"point"`
	List  []*HistoryItem    `json:"list"` // 升序按时间
}

// ─────────────────────────────────────────────────────────────
// 图表数据（ECharts）
// ─────────────────────────────────────────────────────────────

// ChartSeries ECharts series 数据项
type ChartSeries struct {
	Name    string     `json:"name"`    // 序列名称
	Field   string     `json:"field"`   // 实际匹配的 PLC 字段名
	PointId int        `json:"pointId"` // 点位 ID
	Type    string     `json:"type"`    // ECharts series type
	Data    []*float64 `json:"data"`    // 与 xAxis 一一对应，缺失为 null
	Unit    string     `json:"unit"`    // 单位（tooltip 用）
}

// ChartMeta 图表元信息
type ChartMeta struct {
	DeviceId   int    `json:"deviceId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Interval   string `json:"interval"`
	PointCount int    `json:"pointCount"`
	MaxPoints  int    `json:"maxPoints"`
}

// ChartModel ECharts 通用返回结构
type ChartModel struct {
	Meta   *ChartMeta     `json:"meta"`
	Legend []string       `json:"legend"`
	XAxis  []string       `json:"xAxis"`  // 时间轴标签
	Series []*ChartSeries `json:"series"` // 折线列表
}

// ChartTemperatureReq 温度折线图：设备温度 / 回油温度 / 油箱温度
type ChartTemperatureReq struct {
	g.Meta       `path:"/plc/chart/temperature" method:"get" tags:"PLC前台" summary:"设备温度折线图（设备温度/回油温度/油箱温度）"`
	DeviceId     int    `p:"deviceId"     v:"required|min:1#设备ID不能为空" dc:"设备ID"`
	StartTime    string `p:"startTime"    dc:"开始时间，默认最近1小时"`
	EndTime      string `p:"endTime"      dc:"结束时间，默认现在"`
	FieldTemp    string `p:"fieldTemp"    dc:"设备温度字段名，默认 DEVICE_TEMP"`
	FieldOilBack string `p:"fieldOilBack" dc:"回油温度字段名，默认 OIL_BACK_TEMP"`
	FieldOilTank string `p:"fieldOilTank" dc:"油箱温度字段名，默认 OIL_TANK_TEMP"`
}

type ChartTemperatureRes struct {
	*ChartModel
}

// ChartCurrentReq 电流折线图：单台设备电流
type ChartCurrentReq struct {
	g.Meta       `path:"/plc/chart/current" method:"get" tags:"PLC前台" summary:"单设备电流折线图"`
	DeviceId     int    `p:"deviceId"     v:"required|min:1#设备ID不能为空" dc:"设备ID"`
	StartTime    string `p:"startTime"    dc:"开始时间，默认最近1小时"`
	EndTime      string `p:"endTime"      dc:"结束时间，默认现在"`
	FieldCurrent string `p:"fieldCurrent" dc:"电流字段名，默认 CURRENT"`
}

type ChartCurrentRes struct {
	*ChartModel
}
