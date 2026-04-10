// Package service PLC 服务接口定义（手动维护）
package service

import (
	"context"

	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/database/gdb"
)

type (
	// IPlcDevice PLC 设备管理
	IPlcDevice interface {
		Model(ctx context.Context, option ...*handler.Option) *gdb.Model
		List(ctx context.Context, in *sysin.PlcDeviceListInp) (list []*sysin.PlcDeviceListModel, totalCount int, err error)
		View(ctx context.Context, in *sysin.PlcDeviceViewInp) (res *sysin.PlcDeviceViewModel, err error)
		Edit(ctx context.Context, in *sysin.PlcDeviceEditInp) (err error)
		Delete(ctx context.Context, in *sysin.PlcDeviceDeleteInp) (err error)
		Status(ctx context.Context, in *sysin.PlcDeviceStatusInp) (err error)
		ActiveDevices(ctx context.Context) (list []*entity.PlcDevice, err error)
		GetById(ctx context.Context, id int) (dev *entity.PlcDevice, err error)
	}

	// IPlcPoint 数据点管理
	IPlcPoint interface {
		Model(ctx context.Context, option ...*handler.Option) *gdb.Model
		List(ctx context.Context, in *sysin.PlcPointListInp) (list []*sysin.PlcPointListModel, totalCount int, err error)
		View(ctx context.Context, in *sysin.PlcPointViewInp) (res *sysin.PlcPointViewModel, err error)
		Edit(ctx context.Context, in *sysin.PlcPointEditInp) (err error)
		Delete(ctx context.Context, in *sysin.PlcPointDeleteInp) (err error)
		Status(ctx context.Context, in *sysin.PlcPointStatusInp) (err error)
		ActivePoints(ctx context.Context, deviceId int) (list []*entity.PlcPoint, err error)
	}

	// IPlcRealtime 实时数据
	IPlcRealtime interface {
		Get(ctx context.Context, in *sysin.PlcRealtimeInp) (res *sysin.PlcRealtimeModel, err error)
		Set(ctx context.Context, deviceId int, items []*sysin.PlcRealtimeItem) error
	}

	// IPlcHistory 历史记录
	IPlcHistory interface {
		List(ctx context.Context, in *sysin.PlcHistoryInp) (list []*sysin.PlcHistoryModel, totalCount int, err error)
		BatchInsert(ctx context.Context, records []*entity.PlcRecord) error
	}

	// IPlcAlarm 报警管理
	IPlcAlarm interface {
		List(ctx context.Context, in *sysin.PlcAlarmListInp) (list []*sysin.PlcAlarmListModel, totalCount int, err error)
		Resolve(ctx context.Context, in *sysin.PlcAlarmResolveInp) (err error)
		TriggerIfNeeded(ctx context.Context, deviceId int, results []sysin.PlcRealtimeItem) error
	}

	// IPlcApp PLC API应用密钥管理
	IPlcApp interface {
		GetSecretByAppId(ctx context.Context, appId string) (appSecret string, err error)
	}
)

var (
	localPlcDevice   IPlcDevice
	localPlcPoint    IPlcPoint
	localPlcRealtime IPlcRealtime
	localPlcHistory  IPlcHistory
	localPlcAlarm    IPlcAlarm
	localPlcApp      IPlcApp
)

func PlcDevice() IPlcDevice {
	if localPlcDevice == nil {
		panic("implement not found for interface IPlcDevice, forgot register?")
	}
	return localPlcDevice
}

func RegisterPlcDevice(i IPlcDevice) { localPlcDevice = i }

func PlcPoint() IPlcPoint {
	if localPlcPoint == nil {
		panic("implement not found for interface IPlcPoint, forgot register?")
	}
	return localPlcPoint
}

func RegisterPlcPoint(i IPlcPoint) { localPlcPoint = i }

func PlcRealtime() IPlcRealtime {
	if localPlcRealtime == nil {
		panic("implement not found for interface IPlcRealtime, forgot register?")
	}
	return localPlcRealtime
}

func RegisterPlcRealtime(i IPlcRealtime) { localPlcRealtime = i }

func PlcHistory() IPlcHistory {
	if localPlcHistory == nil {
		panic("implement not found for interface IPlcHistory, forgot register?")
	}
	return localPlcHistory
}

func RegisterPlcHistory(i IPlcHistory) { localPlcHistory = i }

func PlcAlarm() IPlcAlarm {
	if localPlcAlarm == nil {
		panic("implement not found for interface IPlcAlarm, forgot register?")
	}
	return localPlcAlarm
}

func RegisterPlcAlarm(i IPlcAlarm) { localPlcAlarm = i }

func PlcApp() IPlcApp {
	if localPlcApp == nil {
		panic("implement not found for interface IPlcApp, forgot register?")
	}
	return localPlcApp
}

func RegisterPlcApp(i IPlcApp) { localPlcApp = i }
