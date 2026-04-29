// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PlcMineDao is the data access object for the table hg_plc_mine.
type PlcMineDao struct {
	table    string
	group    string
	columns  PlcMineColumns
	handlers []gdb.ModelHandler
}

// PlcMineColumns defines and stores column names for the table hg_plc_mine.
type PlcMineColumns struct {
	Id        string
	Name      string
	Location  string
	Remark    string
	Status    string
	CreatedBy string
	UpdatedBy string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

var plcMineColumns = PlcMineColumns{
	Id:        "id",
	Name:      "name",
	Location:  "location",
	Remark:    "remark",
	Status:    "status",
	CreatedBy: "created_by",
	UpdatedBy: "updated_by",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

func NewPlcMineDao(handlers ...gdb.ModelHandler) *PlcMineDao {
	return &PlcMineDao{
		group:    "default",
		table:    "hg_plc_mine",
		columns:  plcMineColumns,
		handlers: handlers,
	}
}

func (dao *PlcMineDao) DB() gdb.DB              { return g.DB(dao.group) }
func (dao *PlcMineDao) Table() string            { return dao.table }
func (dao *PlcMineDao) Columns() PlcMineColumns  { return dao.columns }
func (dao *PlcMineDao) Group() string            { return dao.group }

func (dao *PlcMineDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *PlcMineDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
