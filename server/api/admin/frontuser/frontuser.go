// Package frontuser 前端用户管理接口
package frontuser

import (
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

type ListReq struct {
	g.Meta `path:"/frontUser/list" method:"get" tags:"前端用户" summary:"获取前端用户列表"`
	sysin.FrontUserListInp
}
type ListRes struct {
	form.PageRes
	List []*sysin.FrontUserListModel `json:"list"`
}

type ViewReq struct {
	g.Meta `path:"/frontUser/view" method:"get" tags:"前端用户" summary:"获取前端用户详情"`
	sysin.FrontUserViewInp
}
type ViewRes struct {
	*sysin.FrontUserViewModel
}

type EditReq struct {
	g.Meta `path:"/frontUser/edit" method:"post" tags:"前端用户" summary:"新增/编辑前端用户"`
	sysin.FrontUserEditInp
}
type EditRes struct{}

type DeleteReq struct {
	g.Meta `path:"/frontUser/delete" method:"post" tags:"前端用户" summary:"删除前端用户"`
	sysin.FrontUserDeleteInp
}
type DeleteRes struct{}

type StatusReq struct {
	g.Meta `path:"/frontUser/status" method:"post" tags:"前端用户" summary:"修改前端用户状态"`
	sysin.FrontUserStatusInp
}
type StatusRes struct{}

type ResetPwdReq struct {
	g.Meta `path:"/frontUser/resetPwd" method:"post" tags:"前端用户" summary:"重置前端用户密码"`
	sysin.FrontUserResetPwdInp
}
type ResetPwdRes struct{}
