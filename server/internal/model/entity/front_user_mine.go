// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// FrontUserMine is the golang structure for table hg_front_user_mine.
type FrontUserMine struct {
	Id        int64       `json:"id"        orm:"id"         description:"主键"`
	UserId    int64       `json:"userId"    orm:"user_id"    description:"前端用户ID"`
	MineId    int         `json:"mineId"    orm:"mine_id"    description:"矿场ID"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
}
