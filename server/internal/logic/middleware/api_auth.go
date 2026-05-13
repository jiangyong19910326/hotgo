// Package middleware
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"hotgo/internal/consts"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/response"
	"hotgo/internal/library/token"
	"hotgo/utility/simple"
)

// ApiAuth API鉴权中间件
func (s *sMiddleware) ApiAuth(r *ghttp.Request) {
	var (
		ctx  = r.Context()
		path = gstr.Replace(r.URL.Path, simple.RouterPrefix(ctx, consts.AppApi), "", 1)
	)

	// 不需要验证登录的路由地址
	if s.IsExceptLogin(ctx, consts.AppApi, path) {
		r.Middleware.Next()
		return
	}

	// /api 只接受前台用户 Token，不能走 AdminSite.BindUserContext，否则 admin token 会触发后台角色查询。
	user, err := token.ParseLoginUser(r)
	if err != nil {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), err.Error())
		return
	}
	if user == nil || user.App != consts.AppApi {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "请使用前台用户Token访问")
		return
	}
	contexts.SetUser(r.Context(), user)

	// 验证路由访问权限
	// ...

	r.Middleware.Next()
}
