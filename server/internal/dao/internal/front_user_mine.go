// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FrontUserMineDao is the data access object for the table hg_front_user_mine.
type FrontUserMineDao struct {
	table    string
	group    string
	columns  FrontUserMineColumns
	handlers []gdb.ModelHandler
}

// FrontUserMineColumns defines and stores column names for the table hg_front_user_mine.
type FrontUserMineColumns struct {
	Id        string
	UserId    string
	MineId    string
	CreatedAt string
}

var frontUserMineColumns = FrontUserMineColumns{
	Id:        "id",
	UserId:    "user_id",
	MineId:    "mine_id",
	CreatedAt: "created_at",
}

func NewFrontUserMineDao(handlers ...gdb.ModelHandler) *FrontUserMineDao {
	return &FrontUserMineDao{
		group:    "default",
		table:    "hg_front_user_mine",
		columns:  frontUserMineColumns,
		handlers: handlers,
	}
}

func (dao *FrontUserMineDao) DB() gdb.DB                     { return g.DB(dao.group) }
func (dao *FrontUserMineDao) Table() string                   { return dao.table }
func (dao *FrontUserMineDao) Columns() FrontUserMineColumns   { return dao.columns }
func (dao *FrontUserMineDao) Group() string                   { return dao.group }

func (dao *FrontUserMineDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *FrontUserMineDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
