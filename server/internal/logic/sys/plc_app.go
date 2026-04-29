// Package sys PLC API 应用密钥 — 供签名中间件查询 + 后台管理
package sys

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"hotgo/internal/dao"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

type sPlcApp struct{}

func init() {
	service.RegisterPlcApp(newPlcApp())
}

func newPlcApp() *sPlcApp { return &sPlcApp{} }

const plcAppCacheKey = "plc:app:%s" // %s = appId

func (s *sPlcApp) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.PlcApp.Ctx(ctx), option...)
}

// invalidateCache 清除指定 appId 的 Redis 缓存
func (s *sPlcApp) invalidateCache(ctx context.Context, appId string) {
	if appId == "" {
		return
	}
	_, _ = g.Redis().Del(ctx, fmt.Sprintf(plcAppCacheKey, appId))
}

// GetSecretByAppId 查询应用密钥（带 Redis 缓存，60s TTL）
func (s *sPlcApp) GetSecretByAppId(ctx context.Context, appId string) (appSecret string, err error) {
	cacheKey := fmt.Sprintf(plcAppCacheKey, appId)

	cached, err := g.Redis().Get(ctx, cacheKey)
	if err == nil && !cached.IsNil() && !cached.IsEmpty() {
		return cached.String(), nil
	}

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

	_, _ = g.Redis().Set(ctx, cacheKey, app.AppSecret)
	_, _ = g.Redis().Expire(ctx, cacheKey, 60)

	return app.AppSecret, nil
}

// List 列表
func (s *sPlcApp) List(ctx context.Context, in *sysin.PlcAppListInp) (list []*sysin.PlcAppListModel, totalCount int, err error) {
	mod := s.Model(ctx)
	if in.AppId != "" {
		mod = mod.WhereLike(dao.PlcApp.Columns().AppId, "%"+in.AppId+"%")
	}
	if in.Name != "" {
		mod = mod.WhereLike(dao.PlcApp.Columns().Name, "%"+in.Name+"%")
	}
	if in.Status > 0 {
		mod = mod.Where(dao.PlcApp.Columns().Status, in.Status)
	}

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	err = mod.Page(in.Page, in.PerPage).
		OrderDesc(dao.PlcApp.Columns().Id).
		Scan(&list)
	return
}

// View 详情
func (s *sPlcApp) View(ctx context.Context, in *sysin.PlcAppViewInp) (res *sysin.PlcAppViewModel, err error) {
	res = new(sysin.PlcAppViewModel)
	err = s.Model(ctx).Where(dao.PlcApp.Columns().Id, in.Id).Scan(res)
	return
}

// Edit 新增 / 修改
func (s *sPlcApp) Edit(ctx context.Context, in *sysin.PlcAppEditInp) (err error) {
	if in.Status == 0 {
		in.Status = 1
	}

	if in.Id > 0 {
		// 取出旧 appId 用于失效缓存（旧 appId 改了 secret 也得清旧 key）
		var old entity.PlcApp
		_ = dao.PlcApp.Ctx(ctx).Where(dao.PlcApp.Columns().Id, in.Id).Scan(&old)

		_, err = s.Model(ctx).Where(dao.PlcApp.Columns().Id, in.Id).Data(g.Map{
			dao.PlcApp.Columns().AppId:     in.AppId,
			dao.PlcApp.Columns().AppSecret: in.AppSecret,
			dao.PlcApp.Columns().Name:      in.Name,
			dao.PlcApp.Columns().Remark:    in.Remark,
			dao.PlcApp.Columns().Status:    in.Status,
			dao.PlcApp.Columns().UpdatedAt: gtime.Now(),
		}).Update()
		if err == nil {
			s.invalidateCache(ctx, old.AppId)
			s.invalidateCache(ctx, in.AppId)
		}
	} else {
		_, err = s.Model(ctx).Data(&entity.PlcApp{
			AppId:     in.AppId,
			AppSecret: in.AppSecret,
			Name:      in.Name,
			Remark:    in.Remark,
			Status:    in.Status,
			CreatedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
		}).Insert()
	}
	return
}

// Delete 删除
func (s *sPlcApp) Delete(ctx context.Context, in *sysin.PlcAppDeleteInp) (err error) {
	var rows []entity.PlcApp
	_ = dao.PlcApp.Ctx(ctx).Where(dao.PlcApp.Columns().Id, in.Id).Scan(&rows)
	_, err = s.Model(ctx).Where(dao.PlcApp.Columns().Id, in.Id).Delete()
	if err == nil {
		for _, r := range rows {
			s.invalidateCache(ctx, r.AppId)
		}
	}
	return
}

// Status 改状态
func (s *sPlcApp) Status(ctx context.Context, in *sysin.PlcAppStatusInp) (err error) {
	var old entity.PlcApp
	_ = dao.PlcApp.Ctx(ctx).Where(dao.PlcApp.Columns().Id, in.Id).Scan(&old)
	_, err = s.Model(ctx).Where(dao.PlcApp.Columns().Id, in.Id).
		Data(g.Map{dao.PlcApp.Columns().Status: in.Status, dao.PlcApp.Columns().UpdatedAt: gtime.Now()}).Update()
	if err == nil {
		s.invalidateCache(ctx, old.AppId)
	}
	return
}

// GenSecret 生成 64 字符 hex 随机串
func (s *sPlcApp) GenSecret(_ context.Context) string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
