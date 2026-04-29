// Package plc Open API — PLC 数据接口（签名验签）
package plc

import "github.com/gogf/gf/v2/frame/g"

// ─── 公共签名参数（所有接口 Query 中必传） ───────────────────────────────────────
// appId     string  应用标识
// timestamp string  Unix 秒级时间戳
// nonce     string  随机字符串（8~16 位）
// sign      string  HMAC-SHA256 签名（大写 Hex）

// RealtimeReq 获取设备实时数据
type RealtimeReq struct {
	g.Meta   `path:"/plc/realtime" method:"get" tags:"PLC开放接口" summary:"获取设备实时数据"`
	DeviceId int `p:"deviceId" v:"required|min:1#设备ID不能为空" dc:"设备ID"`
}

// OverviewReq 看板汇总: 设备 + 数据点 + 实时值
type OverviewReq struct {
	g.Meta   `path:"/plc/overview" method:"get" tags:"PLC开放接口" summary:"看板汇总(设备+数据点+实时值)"`
	DeviceId int `p:"deviceId" v:"required|min:1#设备ID不能为空" dc:"设备ID"`
}

// OverviewDevice 设备信息
type OverviewDevice struct {
	Id       int    `json:"id"`
	MineId   int    `json:"mineId"`
	MineName string `json:"mineName"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Remark   string `json:"remark"`
	Status   int    `json:"status"`
}

// OverviewPoint 数据点 + 实时值
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

// OverviewRes 汇总响应
type OverviewRes struct {
	Device *OverviewDevice  `json:"device"`
	Points []*OverviewPoint `json:"points"`
}

// RealtimePointItem 单个点位实时值
type RealtimePointItem struct {
	PointId   int     `json:"pointId"`
	Field     string  `json:"field"`
	Name      string  `json:"name"`
	EngValue  float64 `json:"engValue"`
	Unit      string  `json:"unit"`
	AlarmType int     `json:"alarmType"` // 0正常 1超上限 2超下限
	AlarmMax  *float64 `json:"alarmMax,omitempty"`
	AlarmMin  *float64 `json:"alarmMin,omitempty"`
}

// RealtimeRes 实时数据响应
type RealtimeRes struct {
	DeviceId int                  `json:"deviceId"`
	Points   []*RealtimePointItem `json:"points"`
}

// HistoryReq 获取数据点历史记录
type HistoryReq struct {
	g.Meta    `path:"/plc/history" method:"get" tags:"PLC开放接口" summary:"获取数据点历史记录"`
	PointId   int    `p:"pointId"  v:"required|min:1#点位ID不能为空" dc:"点位ID"`
	StartTime string `p:"startTime" dc:"开始时间，格式 2006-01-02 15:04:05"`
	EndTime   string `p:"endTime"   dc:"结束时间"`
	Page      int    `p:"page"      d:"1"   dc:"页码"`
	PerPage   int    `p:"perPage"   d:"200" dc:"每页数量，最大500"`
}

// HistoryItem 历史记录单条
type HistoryItem struct {
	CollectedAt string   `json:"collectedAt"`
	EngValue    *float64 `json:"engValue"`
	RawValue    string   `json:"rawValue"`
}

// HistoryRes 历史记录响应
type HistoryRes struct {
	TotalCount int            `json:"totalCount"`
	Page       int            `json:"page"`
	PerPage    int            `json:"perPage"`
	List       []*HistoryItem `json:"list"`
}

// AlarmListReq 获取报警列表
type AlarmListReq struct {
	g.Meta     `path:"/plc/alarm/list" method:"get" tags:"PLC开放接口" summary:"获取报警记录列表"`
	DeviceId   int    `p:"deviceId"  dc:"设备ID，不传则查全部"`
	IsResolved int    `p:"isResolved" dc:"处理状态：1已处理 2未处理，不传则查全部"`
	StartTime  string `p:"startTime"  dc:"触发时间起"`
	EndTime    string `p:"endTime"    dc:"触发时间止"`
	Page       int    `p:"page"       d:"1"  dc:"页码"`
	PerPage    int    `p:"perPage"    d:"20" dc:"每页数量"`
}

// AlarmItem 报警记录单条
type AlarmItem struct {
	Id          int      `json:"id"`
	DeviceId    int      `json:"deviceId"`
	PointId     int      `json:"pointId"`
	PointName   string   `json:"pointName"`
	EngValue    float64  `json:"engValue"`
	AlarmType   int      `json:"alarmType"`
	AlarmMin    *float64 `json:"alarmMin,omitempty"`
	AlarmMax    *float64 `json:"alarmMax,omitempty"`
	Unit        string   `json:"unit"`
	IsResolved  int      `json:"isResolved"`
	ResolvedAt  string   `json:"resolvedAt,omitempty"`
	Remark      string   `json:"remark"`
	TriggeredAt string   `json:"triggeredAt"`
}

// AlarmListRes 报警列表响应
type AlarmListRes struct {
	TotalCount int          `json:"totalCount"`
	Page       int          `json:"page"`
	PerPage    int          `json:"perPage"`
	List       []*AlarmItem `json:"list"`
}
