// Package api 前端站点：登录/退出/资料
package api

import (
	"context"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/token"
	"hotgo/internal/model"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type sFrontSite struct{}

func NewFrontSite() *sFrontSite { return &sFrontSite{} }

func init() {
	service.RegisterFrontSite(NewFrontSite())
}

// Login 前端账号密码登录
func (s *sFrontSite) Login(ctx context.Context, in *sysin.FrontUserLoginInp) (res *sysin.FrontUserLoginModel, err error) {
	user, err := service.FrontUser().GetByUsername(ctx, in.Username)
	if err != nil {
		return nil, gerror.Wrap(err, consts.ErrorORM)
	}
	if user == nil {
		return nil, gerror.New("用户名或密码错误")
	}

	if user.Salt == "" {
		return nil, gerror.New("用户信息错误")
	}
	if user.PasswordHash != gmd5.MustEncryptString(in.Password+user.Salt) {
		return nil, gerror.New("用户名或密码错误")
	}
	if user.Status != consts.StatusEnabled {
		return nil, gerror.New("账号已被禁用")
	}

	identity := &model.Identity{
		Id:       user.Id,
		Username: user.Username,
		RealName: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Mobile:   user.Mobile,
		App:      consts.AppApi,
		LoginAt:  gtime.Now(),
	}

	lt, expires, err := token.Login(ctx, identity)
	if err != nil {
		return nil, err
	}

	// 更新最近登录信息
	r := ghttp.RequestFromCtx(ctx)
	loginIp := ""
	if r != nil {
		loginIp = r.GetClientIp()
	}
	cols := dao.FrontUser.Columns()
	if _, e := dao.FrontUser.Ctx(ctx).WherePri(user.Id).Data(g.Map{
		cols.LastLoginAt: gtime.Now(),
		cols.LastLoginIp: loginIp,
	}).Update(); e != nil {
		g.Log().Warningf(ctx, "更新前端用户登录信息失败: %v", e)
	}

	res = &sysin.FrontUserLoginModel{
		Id:       user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Token:    lt,
		Expires:  expires,
	}
	return
}

// Logout 退出登录
func (s *sFrontSite) Logout(ctx context.Context) (err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return gerror.New("非法请求")
	}
	return token.Logout(r)
}

// Profile 当前登录用户资料
func (s *sFrontSite) Profile(ctx context.Context) (res *sysin.FrontUserProfileModel, err error) {
	identity := contexts.GetUser(ctx)
	if identity == nil {
		return nil, gerror.New("未登录")
	}

	user, err := service.FrontUser().GetById(ctx, identity.Id)
	if err != nil {
		return nil, gerror.Wrap(err, consts.ErrorORM)
	}
	if user == nil {
		return nil, gerror.New("用户不存在")
	}

	res = &sysin.FrontUserProfileModel{
		Id:       user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Mobile:   user.Mobile,
		Email:    user.Email,
	}
	return
}
