// Package admin 前端用户管理
package admin

import (
	"context"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
)

type sFrontUser struct{}

func NewFrontUser() *sFrontUser { return &sFrontUser{} }

func init() {
	service.RegisterFrontUser(NewFrontUser())
}

func (s *sFrontUser) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.FrontUser.Ctx(ctx), option...)
}

// List 获取前端用户列表
func (s *sFrontUser) List(ctx context.Context, in *sysin.FrontUserListInp) (list []*sysin.FrontUserListModel, totalCount int, err error) {
	cols := dao.FrontUser.Columns()
	mod := s.Model(ctx)

	if in.Username != "" {
		mod = mod.WhereLike(cols.Username, "%"+in.Username+"%")
	}
	if in.Mobile != "" {
		mod = mod.WhereLike(cols.Mobile, "%"+in.Mobile+"%")
	}
	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}

	totalCount, err = mod.Clone().Count()
	if err != nil {
		return
	}
	err = mod.Page(in.Page, in.PerPage).
		OrderDesc(cols.Id).
		Scan(&list)
	return
}

// View 获取前端用户详情
func (s *sFrontUser) View(ctx context.Context, in *sysin.FrontUserViewInp) (res *sysin.FrontUserViewModel, err error) {
	res = new(sysin.FrontUserViewModel)
	err = s.Model(ctx).WherePri(in.Id).Scan(res)
	return
}

// Edit 新增/编辑前端用户
func (s *sFrontUser) Edit(ctx context.Context, in *sysin.FrontUserEditInp) (err error) {
	if in.Status == 0 {
		in.Status = consts.StatusEnabled
	}

	user := contexts.GetUser(ctx)
	var uid int64
	if user != nil {
		uid = user.Id
	}

	cols := dao.FrontUser.Columns()

	// 编辑
	if in.Id > 0 {
		exist, e := s.Model(ctx).Where(cols.Username, in.Username).WhereNot(cols.Id, in.Id).Count()
		if e != nil {
			return gerror.Wrap(e, "校验用户名失败")
		}
		if exist > 0 {
			return gerror.New("用户名已存在")
		}

		data := g.Map{
			cols.Username:  in.Username,
			cols.Nickname:  in.Nickname,
			cols.Avatar:    in.Avatar,
			cols.Mobile:    in.Mobile,
			cols.Email:     in.Email,
			cols.Remark:    in.Remark,
			cols.Status:    in.Status,
			cols.UpdatedBy: uid,
			cols.UpdatedAt: gtime.Now(),
		}

		// 留空则保持原密码
		if in.Password != "" {
			salt, e := s.Model(ctx).Fields(cols.Salt).WherePri(in.Id).Value()
			if e != nil {
				return gerror.Wrap(e, "获取密码盐失败")
			}
			if salt.IsEmpty() {
				return gerror.New("该用户密码盐缺失")
			}
			data[cols.PasswordHash] = gmd5.MustEncryptString(in.Password + salt.String())
		}

		_, err = s.Model(ctx).WherePri(in.Id).Data(data).Update()
		return
	}

	// 新增
	if in.Password == "" {
		return gerror.New("新增用户密码不能为空")
	}
	exist, err := s.Model(ctx).Where(cols.Username, in.Username).Count()
	if err != nil {
		return gerror.Wrap(err, "校验用户名失败")
	}
	if exist > 0 {
		return gerror.New("用户名已存在")
	}

	salt := grand.S(6)
	row := &entity.FrontUser{
		Username:     in.Username,
		Nickname:     in.Nickname,
		PasswordHash: gmd5.MustEncryptString(in.Password + salt),
		Salt:         salt,
		Avatar:       in.Avatar,
		Mobile:       in.Mobile,
		Email:        in.Email,
		Remark:       in.Remark,
		Status:       in.Status,
		CreatedBy:    uid,
		UpdatedBy:    uid,
		CreatedAt:    gtime.Now(),
		UpdatedAt:    gtime.Now(),
	}
	_, err = s.Model(ctx).Data(row).Insert()
	return
}

// Delete 删除前端用户（软删）
func (s *sFrontUser) Delete(ctx context.Context, in *sysin.FrontUserDeleteInp) (err error) {
	_, err = s.Model(ctx).Where(dao.FrontUser.Columns().Id, in.Id).
		Data(g.Map{dao.FrontUser.Columns().DeletedAt: gtime.Now()}).Update()
	return
}

// Status 更新状态
func (s *sFrontUser) Status(ctx context.Context, in *sysin.FrontUserStatusInp) (err error) {
	_, err = s.Model(ctx).WherePri(in.Id).
		Data(g.Map{dao.FrontUser.Columns().Status: in.Status}).Update()
	return
}

// ResetPwd 重置密码
func (s *sFrontUser) ResetPwd(ctx context.Context, in *sysin.FrontUserResetPwdInp) (err error) {
	cols := dao.FrontUser.Columns()
	salt, err := s.Model(ctx).Fields(cols.Salt).WherePri(in.Id).Value()
	if err != nil {
		return gerror.Wrap(err, "获取密码盐失败")
	}
	if salt.IsEmpty() {
		return gerror.New("用户不存在")
	}
	_, err = s.Model(ctx).WherePri(in.Id).Data(g.Map{
		cols.PasswordHash: gmd5.MustEncryptString(in.Password + salt.String()),
		cols.UpdatedAt:    gtime.Now(),
	}).Update()
	return
}

// GetByUsername 用户名查用户（含密码字段，仅内部调用）
func (s *sFrontUser) GetByUsername(ctx context.Context, username string) (user *entity.FrontUser, err error) {
	user = new(entity.FrontUser)
	err = dao.FrontUser.Ctx(ctx).Where(dao.FrontUser.Columns().Username, username).Scan(user)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, nil
	}
	return
}

// GetById 通过 ID 获取
func (s *sFrontUser) GetById(ctx context.Context, id int64) (user *entity.FrontUser, err error) {
	user = new(entity.FrontUser)
	err = dao.FrontUser.Ctx(ctx).WherePri(id).Scan(user)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, nil
	}
	return
}
