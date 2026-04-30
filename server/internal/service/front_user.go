// Package service 前端用户服务接口定义（手动维护）
package service

import (
	"context"

	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/database/gdb"
)

type (
	// IFrontUser 前端用户管理（后台 CRUD）
	IFrontUser interface {
		Model(ctx context.Context, option ...*handler.Option) *gdb.Model
		List(ctx context.Context, in *sysin.FrontUserListInp) (list []*sysin.FrontUserListModel, totalCount int, err error)
		View(ctx context.Context, in *sysin.FrontUserViewInp) (res *sysin.FrontUserViewModel, err error)
		Edit(ctx context.Context, in *sysin.FrontUserEditInp) (err error)
		Delete(ctx context.Context, in *sysin.FrontUserDeleteInp) (err error)
		Status(ctx context.Context, in *sysin.FrontUserStatusInp) (err error)
		ResetPwd(ctx context.Context, in *sysin.FrontUserResetPwdInp) (err error)
		GetByUsername(ctx context.Context, username string) (user *entity.FrontUser, err error)
		GetById(ctx context.Context, id int64) (user *entity.FrontUser, err error)
	}

	// IFrontSite 前端站点（登录/资料/退出）
	IFrontSite interface {
		Login(ctx context.Context, in *sysin.FrontUserLoginInp) (res *sysin.FrontUserLoginModel, err error)
		Logout(ctx context.Context) (err error)
		Profile(ctx context.Context) (res *sysin.FrontUserProfileModel, err error)
	}
)

var (
	localFrontUser IFrontUser
	localFrontSite IFrontSite
)

func FrontUser() IFrontUser {
	if localFrontUser == nil {
		panic("implement not found for interface IFrontUser, forgot register?")
	}
	return localFrontUser
}

func RegisterFrontUser(i IFrontUser) { localFrontUser = i }

func FrontSite() IFrontSite {
	if localFrontSite == nil {
		panic("implement not found for interface IFrontSite, forgot register?")
	}
	return localFrontSite
}

func RegisterFrontSite(i IFrontSite) { localFrontSite = i }
