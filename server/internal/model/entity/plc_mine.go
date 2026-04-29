// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PlcMine is the golang structure for table hg_plc_mine.
type PlcMine struct {
	Id        int         `json:"id"        orm:"id"         description:"主键"`
	Name      string      `json:"name"      orm:"name"       description:"矿场名称"`
	Location  string      `json:"location"  orm:"location"   description:"地理位置"`
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`
	Status    int         `json:"status"    orm:"status"     description:"状态：1启用 2禁用"`
	CreatedBy int64       `json:"createdBy" orm:"created_by" description:"创建者"`
	UpdatedBy int64       `json:"updatedBy" orm:"updated_by" description:"更新者"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"修改时间"`
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"删除时间"`
}
