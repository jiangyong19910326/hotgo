// Package sys PLC 报警管理
package sys

import (
	"context"

	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sPlcAlarm struct{}

func NewPlcAlarm() *sPlcAlarm { return &sPlcAlarm{} }

func init() {
	service.RegisterPlcAlarm(NewPlcAlarm())
}

func (s *sPlcAlarm) List(ctx context.Context, in *sysin.PlcAlarmListInp) (list []*sysin.PlcAlarmListModel, totalCount int, err error) {
	mod := dao.PlcAlarm.Ctx(ctx)
	if in.DeviceId > 0 {
		mod = mod.Where(dao.PlcAlarm.Columns().DeviceId, in.DeviceId)
	}
	if in.IsResolved > 0 {
		mod = mod.Where(dao.PlcAlarm.Columns().IsResolved, in.IsResolved)
	}
	if in.StartTime != "" {
		mod = mod.WhereGTE(dao.PlcAlarm.Columns().TriggeredAt, in.StartTime)
	}
	if in.EndTime != "" {
		mod = mod.WhereLTE(dao.PlcAlarm.Columns().TriggeredAt, in.EndTime)
	}

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	err = mod.Page(in.Page, in.PerPage).
		OrderDesc(dao.PlcAlarm.Columns().TriggeredAt).
		Scan(&list)
	return
}

// Resolve 处理报警
func (s *sPlcAlarm) Resolve(ctx context.Context, in *sysin.PlcAlarmResolveInp) (err error) {
	user := contexts.GetUser(ctx)
	var uid *int64
	if user != nil {
		id := user.Id
		uid = &id
	}
	now := gtime.Now()
	_, err = dao.PlcAlarm.Ctx(ctx).
		Where(dao.PlcAlarm.Columns().Id, in.Id).
		Data(g.Map{
			dao.PlcAlarm.Columns().IsResolved: 1,
			dao.PlcAlarm.Columns().ResolvedAt: now,
			dao.PlcAlarm.Columns().ResolvedBy: uid,
			dao.PlcAlarm.Columns().Remark:     in.Remark,
		}).Update()
	return
}

// TriggerIfNeeded 根据采集结果同步报警：
// 1. AlarmType != 0 时若没有未处理报警则新增一条；
// 2. AlarmType == 0 且存在未处理报警时自动标记为已恢复。
func (s *sPlcAlarm) TriggerIfNeeded(ctx context.Context, deviceId int, results []sysin.PlcRealtimeItem) error {
	for _, r := range results {
		if r.AlarmType != 0 {
			count, err := dao.PlcAlarm.Ctx(ctx).
				Where(dao.PlcAlarm.Columns().PointId, r.PointId).
				Where(dao.PlcAlarm.Columns().IsResolved, 2).
				Count()
			if err != nil {
				return err
			}
			if count > 0 {
				continue
			}
			now := gtime.Now()
			if _, err = dao.PlcAlarm.Ctx(ctx).Data(&entity.PlcAlarm{
				DeviceId:    deviceId,
				PointId:     r.PointId,
				Field:       r.Field,
				PointName:   r.Name,
				EngValue:    r.EngValue,
				AlarmType:   r.AlarmType,
				AlarmMin:    r.AlarmMin,
				AlarmMax:    r.AlarmMax,
				Unit:        r.Unit,
				IsResolved:  2,
				TriggeredAt: now,
				CreatedAt:   now,
			}).Insert(); err != nil {
				return err
			}
			continue
		}
		// AlarmType == 0：自动恢复该点位仍未处理的报警
		now := gtime.Now()
		_, err := dao.PlcAlarm.Ctx(ctx).
			Where(dao.PlcAlarm.Columns().PointId, r.PointId).
			Where(dao.PlcAlarm.Columns().IsResolved, 2).
			Data(g.Map{
				dao.PlcAlarm.Columns().IsResolved: 1,
				dao.PlcAlarm.Columns().ResolvedAt: now,
				dao.PlcAlarm.Columns().Remark:     "自动恢复",
			}).Update()
		if err != nil {
			return err
		}
	}
	return nil
}
