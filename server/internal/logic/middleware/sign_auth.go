// Package middleware HMAC-SHA256 签名验证中间件
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"

	"hotgo/internal/library/response"
	"hotgo/internal/service"
)

// SignAuth HMAC-SHA256 签名验证中间件
//
// 客户端每次请求须在 Query 中携带：
//   appId     — 应用标识
//   timestamp — Unix 秒级时间戳，与服务器时差须 ≤ 300 秒
//   nonce     — 随机字符串（建议 8~16 位），防重放
//   sign      — 签名值（大写 Hex），生成方式见下
//
// 签名算法：
//   1. 拼接字符串：appId={appId}&nonce={nonce}&timestamp={timestamp}
//   2. sign = HMAC-SHA256(拼接字符串, appSecret)  → 转为大写十六进制
func (s *sMiddleware) SignAuth(r *ghttp.Request) {
	var (
		ctx       = r.Context()
		appId     = r.Get("appId").String()
		timestamp = r.Get("timestamp").String()
		nonce     = r.Get("nonce").String()
		sign      = r.Get("sign").String()
	)

	if appId == "" || timestamp == "" || nonce == "" || sign == "" {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "缺少签名参数 appId/timestamp/nonce/sign")
		return
	}

	// 时间戳校验（允许前后 5 分钟偏差）
	ts := gconv.Int64(timestamp)
	diff := time.Now().Unix() - ts
	if diff > 300 || diff < -300 {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "时间戳已过期，请同步服务器时间")
		return
	}

	// 查询 AppSecret
	appSecret, err := service.PlcApp().GetSecretByAppId(ctx, appId)
	if err != nil {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "无效的 appId: "+err.Error())
		return
	}

	// 计算期望签名
	signStr := "appId=" + appId + "&nonce=" + nonce + "&timestamp=" + timestamp
	expected := calcHMACSHA256(signStr, appSecret)

	if !strings.EqualFold(expected, sign) {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "签名验证失败")
		return
	}

	r.Middleware.Next()
}

// calcHMACSHA256 返回大写 Hex 字符串
func calcHMACSHA256(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}
