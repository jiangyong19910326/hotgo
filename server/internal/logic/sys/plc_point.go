// Package sys PLC 数据点管理
package sys

import (
	"context"

	"hotgo/internal/dao"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sPlcPoint struct{}

func NewPlcPoint() *sPlcPoint { return &sPlcPoint{} }

func init() {
	service.RegisterPlcPoint(NewPlcPoint())
}

func (s *sPlcPoint) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.PlcPoint.Ctx(ctx), option...)
}

func (s *sPlcPoint) List(ctx context.Context, in *sysin.PlcPointListInp) (list []*sysin.PlcPointListModel, totalCount int, err error) {
	mod := s.Model(ctx)
	if in.DeviceId > 0 {
		mod = mod.Where(dao.PlcPoint.Columns().DeviceId, in.DeviceId)
	}
	if in.Name != "" {
		mod = mod.WhereLike(dao.PlcPoint.Columns().Name, "%"+in.Name+"%")
	}
	if in.Status > 0 {
		mod = mod.Where(dao.PlcPoint.Columns().Status, in.Status)
	}

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	err = mod.Page(in.Page, in.PerPage).
		OrderAsc(dao.PlcPoint.Columns().Sort).
		OrderAsc(dao.PlcPoint.Columns().Id).
		Scan(&list)
	return
}

func (s *sPlcPoint) View(ctx context.Context, in *sysin.PlcPointViewInp) (res *sysin.PlcPointViewModel, err error) {
	res = new(sysin.PlcPointViewModel)
	err = s.Model(ctx).Where(dao.PlcPoint.Columns().Id, in.Id).Scan(res)
	return
}

func (s *sPlcPoint) Edit(ctx context.Context, in *sysin.PlcPointEditInp) (err error) {
	if in.Scale == 0 {
		in.Scale = 1
	}
	if in.Status == 0 {
		in.Status = 1
	}

	if in.Id > 0 {
		_, err = s.Model(ctx).Where(dao.PlcPoint.Columns().Id, in.Id).Data(g.Map{
			dao.PlcPoint.Columns().DeviceId:  in.DeviceId,
			dao.PlcPoint.Columns().Name:      in.Name,
			dao.PlcPoint.Columns().Field:     in.Field,
			dao.PlcPoint.Columns().DataType:  in.DataType,
			dao.PlcPoint.Columns().Scale:     in.Scale,
			dao.PlcPoint.Columns().OffsetVal: in.OffsetVal,
			dao.PlcPoint.Columns().Unit:      in.Unit,
			dao.PlcPoint.Columns().AlarmMin:  in.AlarmMin,
			dao.PlcPoint.Columns().AlarmMax:  in.AlarmMax,
			dao.PlcPoint.Columns().Remark:    in.Remark,
			dao.PlcPoint.Columns().Sort:      in.Sort,
			dao.PlcPoint.Columns().Status:    in.Status,
			dao.PlcPoint.Columns().UpdatedAt: gtime.Now(),
		}).Update()
	} else {
		_, err = s.Model(ctx).Data(&entity.PlcPoint{
			DeviceId:  in.DeviceId,
			Name:      in.Name,
			Field:     in.Field,
			DataType:  in.DataType,
			Scale:     in.Scale,
			OffsetVal: in.OffsetVal,
			Unit:      in.Unit,
			AlarmMin:  in.AlarmMin,
			AlarmMax:  in.AlarmMax,
			Remark:    in.Remark,
			Sort:      in.Sort,
			Status:    in.Status,
			CreatedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
		}).Insert()
	}
	return
}

func (s *sPlcPoint) Delete(ctx context.Context, in *sysin.PlcPointDeleteInp) (err error) {
	_, err = s.Model(ctx).Where(dao.PlcPoint.Columns().Id, in.Id).Delete()
	return
}

func (s *sPlcPoint) Status(ctx context.Context, in *sysin.PlcPointStatusInp) (err error) {
	_, err = s.Model(ctx).Where(dao.PlcPoint.Columns().Id, in.Id).
		Data(g.Map{dao.PlcPoint.Columns().Status: in.Status}).Update()
	return
}

// ActivePoints 获取指定设备的所有启用数据点（供 MQTT 订阅器使用）
func (s *sPlcPoint) ActivePoints(ctx context.Context, deviceId int) (list []*entity.PlcPoint, err error) {
	err = s.Model(ctx).
		Where(dao.PlcPoint.Columns().DeviceId, deviceId).
		Where(dao.PlcPoint.Columns().Status, 1).
		OrderAsc(dao.PlcPoint.Columns().Sort).
		Scan(&list)
	return
}

// CreateByFields 按 MQTT payload 字段批量创建数据点（自动接入用）
// 已存在 (device_id, field) 的会被 unique 索引拦截，跳过；返回最终全量列表。
func (s *sPlcPoint) CreateByFields(ctx context.Context, deviceId int, items []service.PlcPointAuto) (list []*entity.PlcPoint, err error) {
	if len(items) == 0 {
		return s.ActivePoints(ctx, deviceId)
	}
	now := gtime.Now()
	rows := make([]entity.PlcPoint, 0, len(items))
	for i, it := range items {
		dt := it.DataType
		if dt == "" {
			dt = "Real"
		}
		rows = append(rows, entity.PlcPoint{
			DeviceId:  deviceId,
			Name:      it.Field,
			Field:     it.Field,
			DataType:  dt,
			Scale:     1,
			OffsetVal: 0,
			Remark:    "MQTT 自动创建",
			Sort:      1000 + i,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	// 逐条插入，忽略 unique 冲突（并发时同字段可能被其他协程先建）
	for _, r := range rows {
		if _, insErr := s.Model(ctx).Data(r).Insert(); insErr != nil {
			// 唯一键冲突跳过
			continue
		}
	}
	return s.ActivePoints(ctx, deviceId)
}
