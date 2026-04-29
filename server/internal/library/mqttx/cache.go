// Package mqttx 设备/点位内存缓存 + MQTT 自动接入
package mqttx

import (
	"context"
	"sync"

	"hotgo/internal/model/entity"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type cacheStore struct {
	mu           sync.RWMutex
	deviceByCode map[string]*entity.PlcDevice // host(DTU编号) → device
	pointsByDev  map[int][]*entity.PlcPoint   // deviceId → points
}

var (
	defaultCache = &cacheStore{
		deviceByCode: make(map[string]*entity.PlcDevice),
		pointsByDev:  make(map[int][]*entity.PlcPoint),
	}
	createMu sync.Mutex // 串行化自动建device/point，避免并发重建
)

// Reload 全量重新加载设备和点位
func Reload(ctx context.Context) error {
	devices, err := service.PlcDevice().ActiveDevices(ctx)
	if err != nil {
		return err
	}

	devMap := make(map[string]*entity.PlcDevice, len(devices))
	pointMap := make(map[int][]*entity.PlcPoint, len(devices))

	for _, dev := range devices {
		if dev.Host == "" {
			continue
		}
		devMap[dev.Host] = dev
		points, perr := service.PlcPoint().ActivePoints(ctx, dev.Id)
		if perr != nil {
			g.Log().Warningf(ctx, "mqttx cache: load points for device %d err: %v", dev.Id, perr)
			continue
		}
		pointMap[dev.Id] = points
	}

	defaultCache.mu.Lock()
	defaultCache.deviceByCode = devMap
	defaultCache.pointsByDev = pointMap
	defaultCache.mu.Unlock()

	g.Log().Infof(ctx, "mqttx cache reloaded: %d devices", len(devMap))
	return nil
}

// GetDeviceByCode 按 DTU 编号查设备
func GetDeviceByCode(code string) *entity.PlcDevice {
	defaultCache.mu.RLock()
	defer defaultCache.mu.RUnlock()
	return defaultCache.deviceByCode[code]
}

// GetPointsByDevice 按设备 ID 查启用点位
func GetPointsByDevice(deviceId int) []*entity.PlcPoint {
	defaultCache.mu.RLock()
	defer defaultCache.mu.RUnlock()
	return defaultCache.pointsByDev[deviceId]
}

// GetOrCreateDevice 缓存miss时按DTU编号自动建设备
func GetOrCreateDevice(ctx context.Context, code string) (*entity.PlcDevice, error) {
	if dev := GetDeviceByCode(code); dev != nil {
		return dev, nil
	}
	createMu.Lock()
	defer createMu.Unlock()
	if dev := GetDeviceByCode(code); dev != nil {
		return dev, nil
	}
	dev, err := service.PlcDevice().CreateByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	defaultCache.mu.Lock()
	defaultCache.deviceByCode[code] = dev
	if _, ok := defaultCache.pointsByDev[dev.Id]; !ok {
		defaultCache.pointsByDev[dev.Id] = nil
	}
	defaultCache.mu.Unlock()
	g.Log().Infof(ctx, "mqttx auto-created device: code=%s, id=%d", code, dev.Id)
	return dev, nil
}

// EnsurePoints 按 MQTT payload 字段补建缺失的 point
// fieldVals: payload中field→value, 用于自动判定data_type (0/1视为Bool)
// 返回该设备最新启用点位列表
func EnsurePoints(ctx context.Context, deviceId int, fieldVals map[string]float64) ([]*entity.PlcPoint, error) {
	existing := GetPointsByDevice(deviceId)
	have := make(map[string]bool, len(existing))
	for _, p := range existing {
		have[p.Field] = true
	}

	var missing []service.PlcPointAuto
	for f, v := range fieldVals {
		if have[f] {
			continue
		}
		dt := "Real"
		if v == 0 || v == 1 {
			dt = "Bool"
		}
		missing = append(missing, service.PlcPointAuto{Field: f, DataType: dt})
	}
	if len(missing) == 0 {
		return existing, nil
	}

	createMu.Lock()
	defer createMu.Unlock()
	// 锁后再核对一次, 避免并发重复建
	existing = GetPointsByDevice(deviceId)
	have = make(map[string]bool, len(existing))
	for _, p := range existing {
		have[p.Field] = true
	}
	final := missing[:0]
	for _, m := range missing {
		if !have[m.Field] {
			final = append(final, m)
		}
	}
	if len(final) == 0 {
		return existing, nil
	}

	all, err := service.PlcPoint().CreateByFields(ctx, deviceId, final)
	if err != nil {
		return existing, err
	}
	defaultCache.mu.Lock()
	defaultCache.pointsByDev[deviceId] = all
	defaultCache.mu.Unlock()
	g.Log().Infof(ctx, "mqttx auto-created %d points for device %d", len(final), deviceId)
	return all, nil
}
