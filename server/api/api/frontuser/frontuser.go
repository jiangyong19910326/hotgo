// Package frontuser 前端登录接口
package frontuser

import (
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

// LoginReq 账号密码登录
type LoginReq struct {
	g.Meta `path:"/frontUser/login" method:"post" tags:"前端用户" summary:"账号登录"`
	sysin.FrontUserLoginInp
}
type LoginRes struct {
	*sysin.FrontUserLoginModel
}

// LogoutReq 退出登录
type LogoutReq struct {
	g.Meta `path:"/frontUser/logout" method:"post" tags:"前端用户" summary:"退出登录"`
}
type LogoutRes struct{}

// ProfileReq 当前登录用户资料
type ProfileReq struct {
	g.Meta `path:"/frontUser/profile" method:"get" tags:"前端用户" summary:"当前登录用户资料"`
}
type ProfileRes struct {
	*sysin.FrontUserProfileModel
}
