// Package sys PLC 设备管理
package sys

import (
	"context"

	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/library/plc"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sPlcDevice struct{}

func NewPlcDevice() *sPlcDevice { return &sPlcDevice{} }

func init() {
	service.RegisterPlcDevice(NewPlcDevice())
}

func (s *sPlcDevice) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.PlcDevice.Ctx(ctx), option...)
}

// List 获取设备列表
func (s *sPlcDevice) List(ctx context.Context, in *sysin.PlcDeviceListInp) (list []*sysin.PlcDeviceListModel, totalCount int, err error) {
	mod := s.Model(ctx)
	if in.Name != "" {
		mod = mod.WhereLike(dao.PlcDevice.Columns().Name, "%"+in.Name+"%")
	}
	if in.Status > 0 {
		mod = mod.Where(dao.PlcDevice.Columns().Status, in.Status)
	}
	mod = mod.WhereNull(dao.PlcDevice.Columns().DeletedAt)

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	err = mod.Page(in.Page, in.PerPage).
		OrderDesc(dao.PlcDevice.Columns().Id).
		Scan(&list)
	return
}

// View 获取设备详情
func (s *sPlcDevice) View(ctx context.Context, in *sysin.PlcDeviceViewInp) (res *sysin.PlcDeviceViewModel, err error) {
	res = new(sysin.PlcDeviceViewModel)
	err = s.Model(ctx).Where(dao.PlcDevice.Columns().Id, in.Id).Scan(res)
	return
}

// Edit 新增/修改设备
func (s *sPlcDevice) Edit(ctx context.Context, in *sysin.PlcDeviceEditInp) (err error) {
	if in.Port <= 0 {
		in.Port = 102
	}
	if in.IntervalMs <= 0 {
		in.IntervalMs = 1000
	}
	if in.Status == 0 {
		in.Status = 1
	}

	user := contexts.GetUser(ctx)
	var uid int64
	if user != nil {
		uid = user.Id
	}

	if in.Id > 0 {
		_, err = s.Model(ctx).Where(dao.PlcDevice.Columns().Id, in.Id).Data(g.Map{
			dao.PlcDevice.Columns().Name:       in.Name,
			dao.PlcDevice.Columns().Host:       in.Host,
			dao.PlcDevice.Columns().Port:       in.Port,
			dao.PlcDevice.Columns().Rack:       in.Rack,
			dao.PlcDevice.Columns().Slot:       in.Slot,
			dao.PlcDevice.Columns().IntervalMs: in.IntervalMs,
			dao.PlcDevice.Columns().Remark:     in.Remark,
			dao.PlcDevice.Columns().Status:     in.Status,
			dao.PlcDevice.Columns().UpdatedBy:  uid,
			dao.PlcDevice.Columns().UpdatedAt:  gtime.Now(),
		}).Update()
		// 设备配置变更，移除旧连接，下次采集时重建
		plc.Remove(in.Id)
	} else {
		_, err = s.Model(ctx).Data(&entity.PlcDevice{
			Name:       in.Name,
			Host:       in.Host,
			Port:       in.Port,
			Rack:       in.Rack,
			Slot:       in.Slot,
			IntervalMs: in.IntervalMs,
			Remark:     in.Remark,
			Status:     in.Status,
			CreatedBy:  uid,
			UpdatedBy:  uid,
			CreatedAt:  gtime.Now(),
			UpdatedAt:  gtime.Now(),
		}).Insert()
	}
	return
}

// Delete 删除设备
func (s *sPlcDevice) Delete(ctx context.Context, in *sysin.PlcDeviceDeleteInp) (err error) {
	mod := s.Model(ctx).Where(dao.PlcDevice.Columns().Id, in.Id)
	_, err = mod.Data(g.Map{
		dao.PlcDevice.Columns().DeletedAt: gtime.Now(),
	}).Update()
	return
}

// Status 更新设备状态
func (s *sPlcDevice) Status(ctx context.Context, in *sysin.PlcDeviceStatusInp) (err error) {
	_, err = s.Model(ctx).Where(dao.PlcDevice.Columns().Id, in.Id).
		Data(g.Map{dao.PlcDevice.Columns().Status: in.Status}).Update()
	if in.Status == 2 {
		plc.Remove(in.Id)
	}
	return
}

// ActiveDevices 获取所有启用中的设备列表（供 Cron 使用）
func (s *sPlcDevice) ActiveDevices(ctx context.Context) (list []*entity.PlcDevice, err error) {
	err = s.Model(ctx).
		Where(dao.PlcDevice.Columns().Status, 1).
		WhereNull(dao.PlcDevice.Columns().DeletedAt).
		Scan(&list)
	return
}

// GetById 通过ID获取设备（内部使用）
func (s *sPlcDevice) GetById(ctx context.Context, id int) (dev *entity.PlcDevice, err error) {
	dev = new(entity.PlcDevice)
	err = s.Model(ctx).Where(dao.PlcDevice.Columns().Id, id).Scan(dev)
	if dev.Id == 0 {
		return nil, gerror.Newf("设备不存在: %d", id)
	}
	return
}
