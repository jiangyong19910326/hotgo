// Package sys 矿场管理
package sys

import (
	"context"

	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sPlcMine struct{}

func NewPlcMine() *sPlcMine { return &sPlcMine{} }

func init() {
	service.RegisterPlcMine(NewPlcMine())
}

func (s *sPlcMine) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.PlcMine.Ctx(ctx), option...)
}

// List 获取矿场列表
func (s *sPlcMine) List(ctx context.Context, in *sysin.PlcMineListInp) (list []*sysin.PlcMineListModel, totalCount int, err error) {
	mod := s.Model(ctx)
	if in.Name != "" {
		mod = mod.WhereLike(dao.PlcMine.Columns().Name, "%"+in.Name+"%")
	}
	if in.Status > 0 {
		mod = mod.Where(dao.PlcMine.Columns().Status, in.Status)
	}
	mod = mod.WhereNull(dao.PlcMine.Columns().DeletedAt)

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	err = mod.Page(in.Page, in.PerPage).
		OrderDesc(dao.PlcMine.Columns().Id).
		Scan(&list)
	return
}

// View 获取矿场详情
func (s *sPlcMine) View(ctx context.Context, in *sysin.PlcMineViewInp) (res *sysin.PlcMineViewModel, err error) {
	res = new(sysin.PlcMineViewModel)
	err = s.Model(ctx).Where(dao.PlcMine.Columns().Id, in.Id).Scan(res)
	return
}

// Edit 新增/修改矿场
func (s *sPlcMine) Edit(ctx context.Context, in *sysin.PlcMineEditInp) (err error) {
	if in.Status == 0 {
		in.Status = 1
	}

	user := contexts.GetUser(ctx)
	var uid int64
	if user != nil {
		uid = user.Id
	}

	if in.Id > 0 {
		_, err = s.Model(ctx).Where(dao.PlcMine.Columns().Id, in.Id).Data(g.Map{
			dao.PlcMine.Columns().Name:      in.Name,
			dao.PlcMine.Columns().Location:  in.Location,
			dao.PlcMine.Columns().Remark:    in.Remark,
			dao.PlcMine.Columns().Status:    in.Status,
			dao.PlcMine.Columns().UpdatedBy: uid,
			dao.PlcMine.Columns().UpdatedAt: gtime.Now(),
		}).Update()
	} else {
		_, err = s.Model(ctx).Data(&entity.PlcMine{
			Name:      in.Name,
			Location:  in.Location,
			Remark:    in.Remark,
			Status:    in.Status,
			CreatedBy: uid,
			UpdatedBy: uid,
			CreatedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
		}).Insert()
	}
	return
}

// Delete 删除矿场
func (s *sPlcMine) Delete(ctx context.Context, in *sysin.PlcMineDeleteInp) (err error) {
	_, err = s.Model(ctx).Where(dao.PlcMine.Columns().Id, in.Id).Data(g.Map{
		dao.PlcMine.Columns().DeletedAt: gtime.Now(),
	}).Update()
	return
}

// Status 更新矿场状态
func (s *sPlcMine) Status(ctx context.Context, in *sysin.PlcMineStatusInp) (err error) {
	_, err = s.Model(ctx).Where(dao.PlcMine.Columns().Id, in.Id).
		Data(g.Map{dao.PlcMine.Columns().Status: in.Status}).Update()
	return
}

// Options 获取启用的矿场下拉列表
func (s *sPlcMine) Options(ctx context.Context) (list []*sysin.PlcMineOption, err error) {
	err = s.Model(ctx).
		Fields(dao.PlcMine.Columns().Id, dao.PlcMine.Columns().Name).
		Where(dao.PlcMine.Columns().Status, 1).
		WhereNull(dao.PlcMine.Columns().DeletedAt).
		OrderAsc(dao.PlcMine.Columns().Id).
		Scan(&list)
	return
}
