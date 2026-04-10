// Package shortlink
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
// @AutoGenerate Version 2.18.6
package shortlink

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

// ListReq 查询短链接列表
type ListReq struct {
	g.Meta `path:"/shortLink/list" method:"get" tags:"短链接" summary:"获取短链接列表"`
	sysin.ShortLinkListInp
}

type ListRes struct {
	form.PageRes
	List []*sysin.ShortLinkListModel `json:"list"   dc:"数据列表"`
}

// ExportReq 导出短链接列表
type ExportReq struct {
	g.Meta `path:"/shortLink/export" method:"get" tags:"短链接" summary:"导出短链接列表"`
	sysin.ShortLinkListInp
}

type ExportRes struct{}

// ViewReq 获取短链接指定信息
type ViewReq struct {
	g.Meta `path:"/shortLink/view" method:"get" tags:"短链接" summary:"获取短链接指定信息"`
	sysin.ShortLinkViewInp
}

type ViewRes struct {
	*sysin.ShortLinkViewModel
}

// EditReq 修改/新增短链接
type EditReq struct {
	g.Meta `path:"/shortLink/edit" method:"post" tags:"短链接" summary:"修改/新增短链接"`
	sysin.ShortLinkEditInp
}

type EditRes struct{}

// DeleteReq 删除短链接
type DeleteReq struct {
	g.Meta `path:"/shortLink/delete" method:"post" tags:"短链接" summary:"删除短链接"`
	sysin.ShortLinkDeleteInp
}

type DeleteRes struct{}

// StatusReq 更新短链接状态
type StatusReq struct {
	g.Meta `path:"/shortLink/status" method:"post" tags:"短链接" summary:"更新短链接状态"`
	sysin.ShortLinkStatusInp
}

type StatusRes struct{}