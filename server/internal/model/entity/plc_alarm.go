// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PlcAlarm is the golang structure for table hg_plc_alarm.
type PlcAlarm struct {
	Id          int64       `json:"id"          orm:"id"           description:"主键"`
	DeviceId    int         `json:"deviceId"    orm:"device_id"    description:"设备 ID"`
	PointId     int         `json:"pointId"     orm:"point_id"     description:"数据点 ID"`
	Field       string      `json:"field"       orm:"field"        description:"字段标识"`
	PointName   string      `json:"pointName"   orm:"point_name"   description:"点位名称"`
	EngValue    float64     `json:"engValue"    orm:"eng_value"    description:"触发时工程值"`
	AlarmType   int         `json:"alarmType"   orm:"alarm_type"   description:"报警类型：1超上限 2超下限"`
	AlarmMin    *float64    `json:"alarmMin"    orm:"alarm_min"    description:"报警下限快照"`
	AlarmMax    *float64    `json:"alarmMax"    orm:"alarm_max"    description:"报警上限快照"`
	Unit        string      `json:"unit"        orm:"unit"         description:"单位"`
	IsResolved  int         `json:"isResolved"  orm:"is_resolved"  description:"是否已处理：1是 2否"`
	ResolvedAt  *gtime.Time `json:"resolvedAt"  orm:"resolved_at"  description:"处理时间"`
	ResolvedBy  *int64      `json:"resolvedBy"  orm:"resolved_by"  description:"处理人 ID"`
	Remark      string      `json:"remark"      orm:"remark"       description:"处理备注"`
	TriggeredAt *gtime.Time `json:"triggeredAt" orm:"triggered_at" description:"触发时间"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`
}
