// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PlcRecordDao is the data access object for the table hg_plc_record.
type PlcRecordDao struct {
	table    string
	group    string
	columns  PlcRecordColumns
	handlers []gdb.ModelHandler
}

// PlcRecordColumns defines and stores column names for the table hg_plc_record.
type PlcRecordColumns struct {
	Id          string
	DeviceId    string
	PointId     string
	Field       string
	RawValue    string
	EngValue    string
	CollectedAt string
}

var plcRecordColumns = PlcRecordColumns{
	Id:          "id",
	DeviceId:    "device_id",
	PointId:     "point_id",
	Field:       "field",
	RawValue:    "raw_value",
	EngValue:    "eng_value",
	CollectedAt: "collected_at",
}

func NewPlcRecordDao(handlers ...gdb.ModelHandler) *PlcRecordDao {
	return &PlcRecordDao{
		group:    "default",
		table:    "hg_plc_record",
		columns:  plcRecordColumns,
		handlers: handlers,
	}
}

func (dao *PlcRecordDao) DB() gdb.DB           { return g.DB(dao.group) }
func (dao *PlcRecordDao) Table() string         { return dao.table }
func (dao *PlcRecordDao) Columns() PlcRecordColumns { return dao.columns }
func (dao *PlcRecordDao) Group() string         { return dao.group }

func (dao *PlcRecordDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *PlcRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
