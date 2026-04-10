// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PlcPointDao is the data access object for the table hg_plc_point.
type PlcPointDao struct {
	table    string
	group    string
	columns  PlcPointColumns
	handlers []gdb.ModelHandler
}

// PlcPointColumns defines and stores column names for the table hg_plc_point.
type PlcPointColumns struct {
	Id         string
	DeviceId   string
	Name       string
	Field      string
	Area       string
	DbNumber   string
	ByteOffset string
	BitOffset  string
	DataType   string
	Scale      string
	OffsetVal  string
	Unit       string
	AlarmMin   string
	AlarmMax   string
	Remark     string
	Sort       string
	Status     string
	CreatedAt  string
	UpdatedAt  string
}

var plcPointColumns = PlcPointColumns{
	Id:         "id",
	DeviceId:   "device_id",
	Name:       "name",
	Field:      "field",
	Area:       "area",
	DbNumber:   "db_number",
	ByteOffset: "byte_offset",
	BitOffset:  "bit_offset",
	DataType:   "data_type",
	Scale:      "scale",
	OffsetVal:  "offset_val",
	Unit:       "unit",
	AlarmMin:   "alarm_min",
	AlarmMax:   "alarm_max",
	Remark:     "remark",
	Sort:       "sort",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

func NewPlcPointDao(handlers ...gdb.ModelHandler) *PlcPointDao {
	return &PlcPointDao{
		group:    "default",
		table:    "hg_plc_point",
		columns:  plcPointColumns,
		handlers: handlers,
	}
}

func (dao *PlcPointDao) DB() gdb.DB          { return g.DB(dao.group) }
func (dao *PlcPointDao) Table() string        { return dao.table }
func (dao *PlcPointDao) Columns() PlcPointColumns { return dao.columns }
func (dao *PlcPointDao) Group() string        { return dao.group }

func (dao *PlcPointDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *PlcPointDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
