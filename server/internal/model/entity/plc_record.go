// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PlcRecord is the golang structure for table hg_plc_record.
type PlcRecord struct {
	Id          int64       `json:"id"          orm:"id"           description:"主键"`
	DeviceId    int         `json:"deviceId"    orm:"device_id"    description:"设备 ID"`
	PointId     int         `json:"pointId"     orm:"point_id"     description:"数据点 ID"`
	Field       string      `json:"field"       orm:"field"        description:"字段标识"`
	RawValue    string      `json:"rawValue"    orm:"raw_value"    description:"原始值"`
	EngValue    *float64    `json:"engValue"    orm:"eng_value"    description:"工程值"`
	CollectedAt *gtime.Time `json:"collectedAt" orm:"collected_at" description:"采集时间"`
}
