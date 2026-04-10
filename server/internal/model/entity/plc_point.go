// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PlcPoint is the golang structure for table hg_plc_point.
type PlcPoint struct {
	Id         int         `json:"id"         orm:"id"          description:"主键"`
	DeviceId   int         `json:"deviceId"   orm:"device_id"   description:"所属设备 ID"`
	Name       string      `json:"name"       orm:"name"        description:"点位名称"`
	Field      string      `json:"field"      orm:"field"       description:"字段标识"`
	Area       string      `json:"area"       orm:"area"        description:"存储区：DB/M/I/Q/V"`
	DbNumber   int         `json:"dbNumber"   orm:"db_number"   description:"DB 块号"`
	ByteOffset int         `json:"byteOffset" orm:"byte_offset" description:"字节偏移量"`
	BitOffset  int         `json:"bitOffset"  orm:"bit_offset"  description:"位偏移量"`
	DataType   string      `json:"dataType"   orm:"data_type"   description:"数据类型"`
	Scale      float64     `json:"scale"      orm:"scale"       description:"换算系数"`
	OffsetVal  float64     `json:"offsetVal"  orm:"offset_val"  description:"换算偏移量"`
	Unit       string      `json:"unit"       orm:"unit"        description:"单位"`
	AlarmMin   *float64    `json:"alarmMin"   orm:"alarm_min"   description:"报警下限"`
	AlarmMax   *float64    `json:"alarmMax"   orm:"alarm_max"   description:"报警上限"`
	Remark     string      `json:"remark"     orm:"remark"      description:"备注"`
	Sort       int         `json:"sort"       orm:"sort"        description:"排序"`
	Status     int         `json:"status"     orm:"status"      description:"状态：1启用 2禁用"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"修改时间"`
}
