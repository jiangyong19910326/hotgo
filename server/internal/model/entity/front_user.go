// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FrontUser is the golang structure for table hg_front_user.
type FrontUser struct {
	Id           int64       `json:"id"           orm:"id"            description:"主键"`
	Username     string      `json:"username"     orm:"username"      description:"用户名"`
	Nickname     string      `json:"nickname"     orm:"nickname"      description:"昵称"`
	PasswordHash string      `json:"passwordHash" orm:"password_hash" description:"密码哈希"`
	Salt         string      `json:"salt"         orm:"salt"          description:"密码盐"`
	Avatar       string      `json:"avatar"       orm:"avatar"        description:"头像"`
	Mobile       string      `json:"mobile"       orm:"mobile"        description:"手机号"`
	Email        string      `json:"email"        orm:"email"         description:"邮箱"`
	Remark       string      `json:"remark"       orm:"remark"        description:"备注"`
	Status       int         `json:"status"       orm:"status"        description:"状态：1启用 2禁用"`
	LastLoginAt  *gtime.Time `json:"lastLoginAt"  orm:"last_login_at" description:"最近登录时间"`
	LastLoginIp  string      `json:"lastLoginIp"  orm:"last_login_ip" description:"最近登录IP"`
	CreatedBy    int64       `json:"createdBy"    orm:"created_by"    description:"创建者"`
	UpdatedBy    int64       `json:"updatedBy"    orm:"updated_by"    description:"更新者"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"修改时间"`
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"删除时间"`
}
