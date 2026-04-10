// Package sys PLC API 应用密钥 — 供签名中间件查询
package sys

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/service"
)

type sPlcApp struct{}

func init() {
	service.RegisterPlcApp(newPlcApp())
}

func newPlcApp() *sPlcApp { return &sPlcApp{} }

const plcAppCacheKey = "plc:app:%s" // %s = appId

// GetSecretByAppId 查询应用密钥（带 Redis 缓存，60s TTL）
func (s *sPlcApp) GetSecretByAppId(ctx context.Context, appId string) (appSecret string, err error) {
	cacheKey := fmt.Sprintf(plcAppCacheKey, appId)

	// 先查 Redis
	cached, err := g.Redis().Get(ctx, cacheKey)
	if err == nil && !cached.IsNil() && !cached.IsEmpty() {
		return cached.String(), nil
	}

	// 查数据库
	var app entity.PlcApp
	err = dao.PlcApp.Ctx(ctx).
		Where(dao.PlcApp.Columns().AppId, appId).
		Where(dao.PlcApp.Columns().Status, 1).
		Scan(&app)
	if err != nil {
		return "", err
	}
	if app.AppId == "" {
		return "", gerror.Newf("appId [%s] 不存在或已禁用", appId)
	}

	// 写入 Redis，TTL 60 秒
	_, _ = g.Redis().Set(ctx, cacheKey, app.AppSecret)
	_, _ = g.Redis().Expire(ctx, cacheKey, 60)

	return app.AppSecret, nil
}
