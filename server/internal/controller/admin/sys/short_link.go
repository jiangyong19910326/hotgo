// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
// @AutoGenerate Version 2.18.6
package sys

import (
	"context"
	"hotgo/api/admin/shortlink"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

var (
	ShortLink = cShortLink{}
)

type cShortLink struct{}

// List 查看短链接列表
func (c *cShortLink) List(ctx context.Context, req *shortlink.ListReq) (res *shortlink.ListRes, err error) {
	list, totalCount, err := service.SysShortLink().List(ctx, &req.ShortLinkListInp)
	if err != nil {
		return
	}

	if list == nil {
		list = []*sysin.ShortLinkListModel{}
	}

	res = new(shortlink.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Export 导出短链接列表
func (c *cShortLink) Export(ctx context.Context, req *shortlink.ExportReq) (res *shortlink.ExportRes, err error) {
	err = service.SysShortLink().Export(ctx, &req.ShortLinkListInp)
	return
}

// Edit 更新短链接
func (c *cShortLink) Edit(ctx context.Context, req *shortlink.EditReq) (res *shortlink.EditRes, err error) {
	err = service.SysShortLink().Edit(ctx, &req.ShortLinkEditInp)
	return
}

// View 获取指定短链接信息
func (c *cShortLink) View(ctx context.Context, req *shortlink.ViewReq) (res *shortlink.ViewRes, err error) {
	data, err := service.SysShortLink().View(ctx, &req.ShortLinkViewInp)
	if err != nil {
		return
	}

	res = new(shortlink.ViewRes)
	res.ShortLinkViewModel = data
	return
}

// Delete 删除短链接
func (c *cShortLink) Delete(ctx context.Context, req *shortlink.DeleteReq) (res *shortlink.DeleteRes, err error) {
	err = service.SysShortLink().Delete(ctx, &req.ShortLinkDeleteInp)
	return
}

// Status 更新短链接状态
func (c *cShortLink) Status(ctx context.Context, req *shortlink.StatusReq) (res *shortlink.StatusRes, err error) {
	err = service.SysShortLink().Status(ctx, &req.ShortLinkStatusInp)
	return
}