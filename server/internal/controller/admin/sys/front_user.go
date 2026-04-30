// Package sys 前端用户管理控制器
package sys

import (
	"context"

	"hotgo/api/admin/frontuser"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

var FrontUser = cFrontUser{}

type cFrontUser struct{}

func (c *cFrontUser) List(ctx context.Context, req *frontuser.ListReq) (res *frontuser.ListRes, err error) {
	list, total, err := service.FrontUser().List(ctx, &req.FrontUserListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.FrontUserListModel{}
	}
	res = new(frontuser.ListRes)
	res.List = list
	res.PageRes.Pack(req, total)
	return
}

func (c *cFrontUser) View(ctx context.Context, req *frontuser.ViewReq) (res *frontuser.ViewRes, err error) {
	data, err := service.FrontUser().View(ctx, &req.FrontUserViewInp)
	if err != nil {
		return
	}
	res = new(frontuser.ViewRes)
	res.FrontUserViewModel = data
	return
}

func (c *cFrontUser) Edit(ctx context.Context, req *frontuser.EditReq) (res *frontuser.EditRes, err error) {
	err = service.FrontUser().Edit(ctx, &req.FrontUserEditInp)
	return
}

func (c *cFrontUser) Delete(ctx context.Context, req *frontuser.DeleteReq) (res *frontuser.DeleteRes, err error) {
	err = service.FrontUser().Delete(ctx, &req.FrontUserDeleteInp)
	return
}

func (c *cFrontUser) Status(ctx context.Context, req *frontuser.StatusReq) (res *frontuser.StatusRes, err error) {
	err = service.FrontUser().Status(ctx, &req.FrontUserStatusInp)
	return
}

func (c *cFrontUser) ResetPwd(ctx context.Context, req *frontuser.ResetPwdReq) (res *frontuser.ResetPwdRes, err error) {
	err = service.FrontUser().ResetPwd(ctx, &req.FrontUserResetPwdInp)
	return
}
