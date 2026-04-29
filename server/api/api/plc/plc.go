// Package plc 前台 PLC 接口（暂未启用签名验签）
package plc

import "github.com/gogf/gf/v2/frame/g"

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
	PointId   int    `p:"pointId"  v:"required|min:1#点位ID不能为空" dc:"点位ID"`
	StartTime string `p:"startTime" dc:"开始时间, 格式 2006-01-02 15:04:05, 为空默认 endTime-1h"`
	EndTime   string `p:"endTime"   dc:"结束时间, 为空默认现在"`
	Limit     int    `p:"limit"     d:"500" dc:"返回点数上限, 默认500, 最大5000"`
}

type HistoryPointInfo struct {
	PointId int    `json:"pointId"`
	Field   string `json:"field"`
	Name    string `json:"name"`
	Unit    string `json:"unit"`
}

type HistoryItem struct {
	Time     string   `json:"time"`               // 采集时间 2006-01-02 15:04:05
	EngValue *float64 `json:"engValue"`           // 工程值
}

type HistoryRes struct {
	Point *HistoryPointInfo `json:"point"`
	List  []*HistoryItem    `json:"list"`           // 升序按时间
}
