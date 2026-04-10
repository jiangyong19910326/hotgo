// Package entity
// This is auto-maintained entity file — do not mix with generated entities.
package entity

import "github.com/gogf/gf/v2/os/gtime"

// PlcApp PLC API应用密钥
type PlcApp struct {
	Id        int         `json:"id"        orm:"id"         description:""`
	AppId     string      `json:"appId"     orm:"app_id"     description:"AppID"`
	AppSecret string      `json:"appSecret" orm:"app_secret" description:"AppSecret"`
	Name      string      `json:"name"      orm:"name"       description:"应用名称"`
	Remark    string      `json:"remark"    orm:"remark"     description:"备注"`
	Status    int         `json:"status"    orm:"status"     description:"1启用 2禁用"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
