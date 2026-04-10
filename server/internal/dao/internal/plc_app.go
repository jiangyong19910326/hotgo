// ==========================================================================
// This file is maintained manually for hg_plc_app.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PlcAppDao is the data access object for the table hg_plc_app.
type PlcAppDao struct {
	table    string
	group    string
	columns  PlcAppColumns
	handlers []gdb.ModelHandler
}

// PlcAppColumns defines and stores column names for the table hg_plc_app.
type PlcAppColumns struct {
	Id        string
	AppId     string
	AppSecret string
	Name      string
	Remark    string
	Status    string
	CreatedAt string
	UpdatedAt string
}

var plcAppColumns = PlcAppColumns{
	Id:        "id",
	AppId:     "app_id",
	AppSecret: "app_secret",
	Name:      "name",
	Remark:    "remark",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

func NewPlcAppDao(handlers ...gdb.ModelHandler) *PlcAppDao {
	return &PlcAppDao{
		group:    "default",
		table:    "hg_plc_app",
		columns:  plcAppColumns,
		handlers: handlers,
	}
}

func (dao *PlcAppDao) DB() gdb.DB              { return g.DB(dao.group) }
func (dao *PlcAppDao) Table() string           { return dao.table }
func (dao *PlcAppDao) Columns() PlcAppColumns  { return dao.columns }
func (dao *PlcAppDao) Group() string           { return dao.group }

func (dao *PlcAppDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *PlcAppDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
