// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PlcDeviceDao is the data access object for the table hg_plc_device.
type PlcDeviceDao struct {
	table    string
	group    string
	columns  PlcDeviceColumns
	handlers []gdb.ModelHandler
}

// PlcDeviceColumns defines and stores column names for the table hg_plc_device.
type PlcDeviceColumns struct {
	Id         string
	Name       string
	Host       string
	Port       string
	Rack       string
	Slot       string
	IntervalMs string
	Remark     string
	Status     string
	CreatedBy  string
	UpdatedBy  string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}

var plcDeviceColumns = PlcDeviceColumns{
	Id:         "id",
	Name:       "name",
	Host:       "host",
	Port:       "port",
	Rack:       "rack",
	Slot:       "slot",
	IntervalMs: "interval_ms",
	Remark:     "remark",
	Status:     "status",
	CreatedBy:  "created_by",
	UpdatedBy:  "updated_by",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

func NewPlcDeviceDao(handlers ...gdb.ModelHandler) *PlcDeviceDao {
	return &PlcDeviceDao{
		group:    "default",
		table:    "hg_plc_device",
		columns:  plcDeviceColumns,
		handlers: handlers,
	}
}

func (dao *PlcDeviceDao) DB() gdb.DB           { return g.DB(dao.group) }
func (dao *PlcDeviceDao) Table() string         { return dao.table }
func (dao *PlcDeviceDao) Columns() PlcDeviceColumns { return dao.columns }
func (dao *PlcDeviceDao) Group() string         { return dao.group }

func (dao *PlcDeviceDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *PlcDeviceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
