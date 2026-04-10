// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PlcAlarmDao is the data access object for the table hg_plc_alarm.
type PlcAlarmDao struct {
	table    string
	group    string
	columns  PlcAlarmColumns
	handlers []gdb.ModelHandler
}

// PlcAlarmColumns defines and stores column names for the table hg_plc_alarm.
type PlcAlarmColumns struct {
	Id          string
	DeviceId    string
	PointId     string
	Field       string
	PointName   string
	EngValue    string
	AlarmType   string
	AlarmMin    string
	AlarmMax    string
	Unit        string
	IsResolved  string
	ResolvedAt  string
	ResolvedBy  string
	Remark      string
	TriggeredAt string
	CreatedAt   string
}

var plcAlarmColumns = PlcAlarmColumns{
	Id:          "id",
	DeviceId:    "device_id",
	PointId:     "point_id",
	Field:       "field",
	PointName:   "point_name",
	EngValue:    "eng_value",
	AlarmType:   "alarm_type",
	AlarmMin:    "alarm_min",
	AlarmMax:    "alarm_max",
	Unit:        "unit",
	IsResolved:  "is_resolved",
	ResolvedAt:  "resolved_at",
	ResolvedBy:  "resolved_by",
	Remark:      "remark",
	TriggeredAt: "triggered_at",
	CreatedAt:   "created_at",
}

func NewPlcAlarmDao(handlers ...gdb.ModelHandler) *PlcAlarmDao {
	return &PlcAlarmDao{
		group:    "default",
		table:    "hg_plc_alarm",
		columns:  plcAlarmColumns,
		handlers: handlers,
	}
}

func (dao *PlcAlarmDao) DB() gdb.DB          { return g.DB(dao.group) }
func (dao *PlcAlarmDao) Table() string        { return dao.table }
func (dao *PlcAlarmDao) Columns() PlcAlarmColumns { return dao.columns }
func (dao *PlcAlarmDao) Group() string        { return dao.group }

func (dao *PlcAlarmDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *PlcAlarmDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
