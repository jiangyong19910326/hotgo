// Package service PLC 服务接口定义（手动维护）
package service

import (
	"context"

	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/database/gdb"
)

// PlcPointAuto MQTT 自动建点位的最小参数
type PlcPointAuto struct {
	Field    string
	DataType string
}

type (
	// IPlcMine 矿场管理
	IPlcMine interface {
		Model(ctx context.Context, option ...*handler.Option) *gdb.Model
		List(ctx context.Context, in *sysin.PlcMineListInp) (list []*sysin.PlcMineListModel, totalCount int, err error)
		View(ctx context.Context, in *sysin.PlcMineViewInp) (res *sysin.PlcMineViewModel, err error)
		Edit(ctx context.Context, in *sysin.PlcMineEditInp) (err error)
		Delete(ctx context.Context, in *sysin.PlcMineDeleteInp) (err error)
		Status(ctx context.Context, in *sysin.PlcMineStatusInp) (err error)
		Options(ctx context.Context) (list []*sysin.PlcMineOption, err error)
	}

	// IPlcDevice PLC 设备管理
	IPlcDevice interface {
		Model(ctx context.Context, option ...*handler.Option) *gdb.Model
		List(ctx context.Context, in *sysin.PlcDeviceListInp) (list []*sysin.PlcDeviceListModel, totalCount int, err error)
		View(ctx context.Context, in *sysin.PlcDeviceViewInp) (res *sysin.PlcDeviceViewModel, err error)
		Edit(ctx context.Context, in *sysin.PlcDeviceEditInp) (err error)
		Delete(ctx context.Context, in *sysin.PlcDeviceDeleteInp) (err error)
		Status(ctx context.Context, in *sysin.PlcDeviceStatusInp) (err error)
		Control(ctx context.Context, in *sysin.PlcDeviceControlInp) (res *sysin.PlcDeviceControlModel, err error)
		ActiveDevices(ctx context.Context) (list []*entity.PlcDevice, err error)
		GetById(ctx context.Context, id int) (dev *entity.PlcDevice, err error)
		CreateByCode(ctx context.Context, code string) (dev *entity.PlcDevice, err error)
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
		CreateByFields(ctx context.Context, deviceId int, items []PlcPointAuto) (list []*entity.PlcPoint, err error)
	}

	// IPlcRealtime 实时数据
	IPlcRealtime interface {
		Get(ctx context.Context, in *sysin.PlcRealtimeInp) (res *sysin.PlcRealtimeModel, err error)
		Set(ctx context.Context, deviceId int, items []*sysin.PlcRealtimeItem) error
		Overview(ctx context.Context, in *sysin.PlcOverviewInp) (res *sysin.PlcOverviewModel, err error)
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
		Model(ctx context.Context, option ...*handler.Option) *gdb.Model
		GetSecretByAppId(ctx context.Context, appId string) (appSecret string, err error)
		List(ctx context.Context, in *sysin.PlcAppListInp) (list []*sysin.PlcAppListModel, totalCount int, err error)
		View(ctx context.Context, in *sysin.PlcAppViewInp) (res *sysin.PlcAppViewModel, err error)
		Edit(ctx context.Context, in *sysin.PlcAppEditInp) (err error)
		Delete(ctx context.Context, in *sysin.PlcAppDeleteInp) (err error)
		Status(ctx context.Context, in *sysin.PlcAppStatusInp) (err error)
		GenSecret(ctx context.Context) (secret string)
	}

	// IPlcChart 图表数据（ECharts）
	IPlcChart interface {
		// Temperature 返回单台设备的温度折线图数据（设备温度/回油温度/油箱温度）
		Temperature(ctx context.Context, in *sysin.PlcChartTemperatureInp) (res *sysin.PlcChartModel, err error)
		// Current 返回多台设备的电流折线图数据
		Current(ctx context.Context, in *sysin.PlcChartCurrentInp) (res *sysin.PlcChartModel, err error)
	}
)

var (
	localPlcMine     IPlcMine
	localPlcDevice   IPlcDevice
	localPlcPoint    IPlcPoint
	localPlcRealtime IPlcRealtime
	localPlcHistory  IPlcHistory
	localPlcAlarm    IPlcAlarm
	localPlcApp      IPlcApp
	localPlcChart    IPlcChart
)

func PlcMine() IPlcMine {
	if localPlcMine == nil {
		panic("implement not found for interface IPlcMine, forgot register?")
	}
	return localPlcMine
}

func RegisterPlcMine(i IPlcMine) { localPlcMine = i }

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

func PlcChart() IPlcChart {
	if localPlcChart == nil {
		panic("implement not found for interface IPlcChart, forgot register?")
	}
	return localPlcChart
}

func RegisterPlcChart(i IPlcChart) { localPlcChart = i }
