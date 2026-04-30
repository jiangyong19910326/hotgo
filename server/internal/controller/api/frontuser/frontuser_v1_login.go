package frontuser

import (
	"context"

	v1 "hotgo/api/api/frontuser"
	"hotgo/internal/service"
)

// Login 账号登录
func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	data, err := service.FrontSite().Login(ctx, &req.FrontUserLoginInp)
	if err != nil {
		return nil, err
	}
	res = &v1.LoginRes{FrontUserLoginModel: data}
	return
}

// Logout 退出
func (c *ControllerV1) Logout(ctx context.Context, _ *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	err = service.FrontSite().Logout(ctx)
	return
}

// Profile 当前登录用户资料
func (c *ControllerV1) Profile(ctx context.Context, _ *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	data, err := service.FrontSite().Profile(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.ProfileRes{FrontUserProfileModel: data}
	return
}
