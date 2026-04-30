// Package router
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package router

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/controller/api/frontuser"
	"hotgo/internal/controller/api/member"
	"hotgo/internal/controller/api/pay"
	apiplc "hotgo/internal/controller/api/plc"
	"hotgo/internal/controller/api/shortlink"
	"hotgo/internal/service"
	"hotgo/utility/simple"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Api 前台路由
func Api(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group(simple.RouterPrefix(ctx, consts.AppApi), func(group *ghttp.RouterGroup) {
		group.Bind(
			pay.NewV1(),       // 支付异步通知
			shortlink.NewV1(), // 短链接
			apiplc.NewV1(),    // PLC 前台接口（暂未启用签名验签）
		)
		group.Middleware(service.Middleware().ApiAuth)
		group.Bind(
			member.NewV1(),    // 管理员
			frontuser.NewV1(), // 前端用户（/frontUser/login 在 exceptLogin 中放行）
		)
	})
}
