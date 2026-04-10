// Package router Open API 路由（签名验签）
package router

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	"hotgo/internal/consts"
	openplc "hotgo/internal/controller/open/plc"
	"hotgo/internal/service"
	"hotgo/utility/simple"
)

// Open 开放接口路由（所有请求须通过 HMAC-SHA256 签名验证）
func Open(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group(simple.RouterPrefix(ctx, consts.AppOpen), func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().SignAuth)
		group.Bind(
			openplc.NewV1(), // PLC 数据开放接口
		)
	})
}
