// Package sysin 前端用户输入/输出模型
package sysin

import (
	"hotgo/internal/model/input/form"

	"github.com/gogf/gf/v2/os/gtime"
)

// ─────────────────────────────────────────────────────────────
// 后台 CRUD
// ─────────────────────────────────────────────────────────────

type FrontUserListInp struct {
	form.PageReq
	Username string `json:"username" dc:"用户名"`
	Mobile   string `json:"mobile"   dc:"手机号"`
	Status   int    `json:"status"   dc:"状态"`
}

type FrontUserListModel struct {
	Id          int64       `json:"id"`
	Username    string      `json:"username"`
	Nickname    string      `json:"nickname"`
	Avatar      string      `json:"avatar"`
	Mobile      string      `json:"mobile"`
	Email       string      `json:"email"`
	Remark      string      `json:"remark"`
	Status      int         `json:"status"`
	LastLoginAt *gtime.Time `json:"lastLoginAt"`
	LastLoginIp string      `json:"lastLoginIp"`
	CreatedAt   *gtime.Time `json:"createdAt"`
	MineIds     []int       `json:"mineIds"`
	MineNames   []string    `json:"mineNames"`
}

type FrontUserViewInp struct {
	Id int64 `json:"id" v:"required#ID不能为空"`
}

type FrontUserViewModel = FrontUserListModel

type FrontUserEditInp struct {
	Id       int64  `json:"id"`
	Username string `json:"username" v:"required#用户名不能为空" dc:"用户名"`
	Password string `json:"password"                          dc:"密码（新增必填，编辑留空保持不变）"`
	Nickname string `json:"nickname"                          dc:"昵称"`
	Avatar   string `json:"avatar"                            dc:"头像"`
	Mobile   string `json:"mobile"                            dc:"手机号"`
	Email    string `json:"email"                             dc:"邮箱"`
	Remark   string `json:"remark"                            dc:"备注"`
	Status   int    `json:"status"                            dc:"状态"`
	MineIds  []int  `json:"mineIds"                           dc:"绑定矿场ID列表"`
}

type FrontUserDeleteInp struct {
	Id interface{} `json:"id" v:"required#ID不能为空"`
}

type FrontUserStatusInp struct {
	Id     int64 `json:"id"     v:"required#ID不能为空"`
	Status int   `json:"status" v:"required#状态不能为空"`
}

type FrontUserResetPwdInp struct {
	Id       int64  `json:"id"       v:"required#ID不能为空"`
	Password string `json:"password" v:"required|length:6,32#新密码不能为空|新密码长度6~32位"`
}

// ─────────────────────────────────────────────────────────────
// 前端登录
// ─────────────────────────────────────────────────────────────

type FrontUserLoginInp struct {
	Username string `json:"username" v:"required#用户名不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
}

type FrontUserLoginModel struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Token    string `json:"token"`
	Expires  int64  `json:"expires"`
}

type FrontUserProfileModel struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
}
