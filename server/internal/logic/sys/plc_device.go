// Package sys PLC 设备管理
package sys

import (
	"context"
	"encoding/json"
	"strings"

	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/library/mqttx"
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
	d := dao.PlcDevice.Columns()
	mod := s.Model(ctx).
		LeftJoin("hg_plc_mine m", "m.id = "+dao.PlcDevice.Table()+"."+d.MineId)

	if in.MineId > 0 {
		mod = mod.Where(dao.PlcDevice.Table()+"."+d.MineId, in.MineId)
	}
	if in.Name != "" {
		mod = mod.WhereLike(dao.PlcDevice.Table()+"."+d.Name, "%"+in.Name+"%")
	}
	if in.Status > 0 {
		mod = mod.Where(dao.PlcDevice.Table()+"."+d.Status, in.Status)
	}

	totalCount, err = mod.Clone().Count()
	if err != nil {
		return
	}
	err = mod.Fields(dao.PlcDevice.Table()+".*", "m.name as mine_name").
		Page(in.Page, in.PerPage).
		OrderDesc(dao.PlcDevice.Table() + "." + d.Id).
		Scan(&list)
	return
}

// View 获取设备详情
func (s *sPlcDevice) View(ctx context.Context, in *sysin.PlcDeviceViewInp) (res *sysin.PlcDeviceViewModel, err error) {
	res = new(sysin.PlcDeviceViewModel)
	d := dao.PlcDevice.Columns()
	err = s.Model(ctx).
		LeftJoin("hg_plc_mine m", "m.id = "+dao.PlcDevice.Table()+"."+d.MineId).
		Fields(dao.PlcDevice.Table()+".*", "m.name as mine_name").
		Where(dao.PlcDevice.Table()+"."+d.Id, in.Id).
		Scan(res)
	return
}

// Edit 新增/修改设备
func (s *sPlcDevice) Edit(ctx context.Context, in *sysin.PlcDeviceEditInp) (err error) {
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
			dao.PlcDevice.Columns().MineId:    in.MineId,
			dao.PlcDevice.Columns().Name:      in.Name,
			dao.PlcDevice.Columns().Host:      in.Host,
			dao.PlcDevice.Columns().Remark:    in.Remark,
			dao.PlcDevice.Columns().Status:    in.Status,
			dao.PlcDevice.Columns().UpdatedBy: uid,
			dao.PlcDevice.Columns().UpdatedAt: gtime.Now(),
		}).Update()
	} else {
		_, err = s.Model(ctx).Data(&entity.PlcDevice{
			MineId:    in.MineId,
			Name:      in.Name,
			Host:      in.Host,
			Remark:    in.Remark,
			Status:    in.Status,
			CreatedBy: uid,
			UpdatedBy: uid,
			CreatedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
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
	return
}

// Control 通过 MQTT 向设备发送一键启动/停止命令。
func (s *sPlcDevice) Control(ctx context.Context, in *sysin.PlcDeviceControlInp) (res *sysin.PlcDeviceControlModel, err error) {
	g.Log().Infof(ctx, "plc device control requested: deviceId=%d action=%s", in.DeviceId, in.Action)

	dev, err := s.GetById(ctx, in.DeviceId)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(dev.Host) == "" {
		g.Log().Warningf(ctx, "plc device control failed: empty host, deviceId=%d", in.DeviceId)
		return nil, gerror.New("设备DTU编号为空，无法发送控制命令")
	}

	point, err := s.findControlPoint(ctx, in.DeviceId, in.Action)
	if err != nil {
		return nil, err
	}
	if point == nil || strings.TrimSpace(point.Field) == "" {
		g.Log().Warningf(ctx, "plc device control failed: control point not found, deviceId=%d action=%s", in.DeviceId, in.Action)
		return nil, gerror.New("未找到对应的一键启停点位")
	}

	payload, err := json.Marshal(g.Map{point.Field: 1})
	if err != nil {
		return nil, err
	}

	topic := "/dtu/" + strings.Trim(dev.Host, "/") + "/cmd"
	qos := byte(0)
	g.Log().Infof(ctx, "plc device control publish: deviceId=%d action=%s topic=%s qos=%d payload=%s", in.DeviceId, in.Action, topic, qos, string(payload))
	if err = mqttx.Publish(ctx, topic, payload, qos); err != nil {
		return nil, err
	}

	return &sysin.PlcDeviceControlModel{
		DeviceId:   in.DeviceId,
		DeviceCode: dev.Host,
		Action:     in.Action,
		PointId:    point.Id,
		PointField: point.Field,
		Topic:      topic,
		Payload:    string(payload),
		QoS:        qos,
	}, nil
}

func (s *sPlcDevice) findControlPoint(ctx context.Context, deviceId int, action string) (*entity.PlcPoint, error) {
	points, err := service.PlcPoint().ActivePoints(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	candidates := []string{"一键启动", "启动", "start"}
	if strings.EqualFold(action, "stop") {
		candidates = []string{"一键停止", "停止", "stop"}
	}

	for _, point := range points {
		text := strings.ToLower(point.Field + " " + point.Name + " " + point.Remark)
		for _, keyword := range candidates {
			if strings.Contains(text, strings.ToLower(keyword)) {
				return point, nil
			}
		}
	}

	return nil, nil
}

// ActiveDevices 获取所有启用中的设备列表（供 MQTT 订阅器使用）
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

// CreateByCode 按 DTU 编号创建设备（MQTT 自动接入用）
func (s *sPlcDevice) CreateByCode(ctx context.Context, code string) (dev *entity.PlcDevice, err error) {
	now := gtime.Now()
	row := &entity.PlcDevice{
		MineId:    0,
		Name:      code,
		Host:      code,
		Remark:    "MQTT 自动接入",
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	res, err := s.Model(ctx).Data(row).Insert()
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	row.Id = int(id)
	return row, nil
}
