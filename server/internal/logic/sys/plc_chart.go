// Package sys PLC 图表数据（ECharts）
package sys

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"hotgo/internal/dao"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

const plcChartMaxPoints = 1000

type sPlcChart struct{}

type plcChartPoint struct {
	Id    int    `orm:"id"`
	Field string `orm:"field"`
	Unit  string `orm:"unit"`
}

type plcChartRange struct {
	Start    string
	End      string
	Interval string
}

func NewPlcChart() *sPlcChart { return &sPlcChart{} }

func init() {
	service.RegisterPlcChart(NewPlcChart())
}

// defaultChartRange 时间窗策略：
// 1. 显式传 startTime/endTime 时按用户指定。
// 2. 缺失 endTime 时优先取该设备最新一条记录的时间，回退到 now。
// 3. 缺失 startTime 时按 endTime 向前 1 小时。
func defaultChartRange(ctx context.Context, deviceId int, startTime, endTime string) plcChartRange {
	end := normalizeChartTime(endTime, false)
	if end == "" {
		if latest := latestRecordTime(ctx, deviceId); latest != "" {
			end = latest
		} else {
			end = gtime.Now().Format("Y-m-d H:i:s")
		}
	}
	start := normalizeChartTime(startTime, true)
	if start == "" {
		if et, err := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local); err == nil {
			start = et.Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
		} else {
			start = gtime.NewFromTime(time.Now().Add(-1 * time.Hour)).Format("Y-m-d H:i:s")
		}
	}
	return plcChartRange{Start: start, End: end, Interval: chooseChartInterval(start, end)}
}

// latestRecordTime 返回该设备最新一条 plc_record 的采集时间字符串；无记录时返回空串。
func latestRecordTime(ctx context.Context, deviceId int) string {
	if deviceId <= 0 {
		return ""
	}
	var row struct {
		CollectedAt *gtime.Time `orm:"collected_at"`
	}
	_ = dao.PlcRecord.Ctx(ctx).
		Fields(dao.PlcRecord.Columns().CollectedAt).
		Where(dao.PlcRecord.Columns().DeviceId, deviceId).
		OrderDesc(dao.PlcRecord.Columns().CollectedAt).
		Limit(1).
		Scan(&row)
	if row.CollectedAt == nil {
		return ""
	}
	return row.CollectedAt.Format("Y-m-d H:i:s")
}

// normalizeChartTime 兼容前端传 YYYY-MM、YYYY-MM-DD、YYYY-MM-DD HH:mm:ss 三种格式。
func normalizeChartTime(value string, isStart bool) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	switch len(value) {
	case len("2006-01"):
		if t, err := time.ParseInLocation("2006-01", value, time.Local); err == nil {
			if isStart {
				return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
			}
			return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 0, time.Local).Format("2006-01-02 15:04:05")
		}
	case len("2006-01-02"):
		if isStart {
			return value + " 00:00:00"
		}
		return value + " 23:59:59"
	}
	return value
}

func chooseChartInterval(start, end string) string {
	st, err1 := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	et, err2 := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	if err1 != nil || err2 != nil || !et.After(st) {
		return "1m"
	}
	d := et.Sub(st)
	switch {
	case d <= 6*time.Hour:
		return "1m"
	case d <= 3*24*time.Hour:
		return "5m"
	case d <= 30*24*time.Hour:
		return "1h"
	default:
		return "1d"
	}
}

func bucketExpr(interval string) string {
	switch interval {
	case "5m":
		return "DATE_FORMAT(FROM_UNIXTIME(FLOOR(UNIX_TIMESTAMP(collected_at)/300)*300), '%Y-%m-%d %H:%i:00')"
	case "1h":
		return "DATE_FORMAT(collected_at, '%Y-%m-%d %H:00:00')"
	case "1d":
		return "DATE_FORMAT(collected_at, '%Y-%m-%d 00:00:00')"
	default:
		return "DATE_FORMAT(collected_at, '%Y-%m-%d %H:%i:00')"
	}
}

func (s *sPlcChart) findPoint(ctx context.Context, deviceId int, field string) (*plcChartPoint, error) {
	if field == "" {
		return nil, nil
	}
	var point plcChartPoint
	err := dao.PlcPoint.Ctx(ctx).
		Fields(dao.PlcPoint.Columns().Id, dao.PlcPoint.Columns().Field, dao.PlcPoint.Columns().Unit).
		Where(dao.PlcPoint.Columns().DeviceId, deviceId).
		Where(dao.PlcPoint.Columns().Field, field).
		Where(dao.PlcPoint.Columns().Status, 1).
		Scan(&point)
	if err != nil || point.Id == 0 {
		return nil, err
	}
	return &point, nil
}

func (s *sPlcChart) queryPointData(ctx context.Context, pointId int, r plcChartRange) (map[string]float64, error) {
	var rows []struct {
		Bucket string   `orm:"bucket"`
		Value  *float64 `orm:"value"`
	}
	err := dao.PlcRecord.Ctx(ctx).
		Fields(fmt.Sprintf("%s AS bucket, AVG(%s) AS value", bucketExpr(r.Interval), dao.PlcRecord.Columns().EngValue)).
		Where(dao.PlcRecord.Columns().PointId, pointId).
		WhereGTE(dao.PlcRecord.Columns().CollectedAt, r.Start).
		WhereLTE(dao.PlcRecord.Columns().CollectedAt, r.End).
		WhereNotNull(dao.PlcRecord.Columns().EngValue).
		Group("bucket").
		OrderAsc("bucket").
		Limit(plcChartMaxPoints).
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	result := make(map[string]float64, len(rows))
	for _, row := range rows {
		if row.Bucket == "" || row.Value == nil {
			continue
		}
		result[row.Bucket] = *row.Value
	}
	return result, nil
}

func collectXAxis(maps ...map[string]float64) []string {
	seen := make(map[string]struct{})
	for _, m := range maps {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	labels := make([]string, 0, len(seen))
	for k := range seen {
		labels = append(labels, k)
	}
	sort.Strings(labels)
	return labels
}

func buildSeries(name string, point *plcChartPoint, defaultUnit string, xAxis []string, data map[string]float64) *sysin.PlcChartSeries {
	values := make([]*float64, len(xAxis))
	for i, label := range xAxis {
		if v, ok := data[label]; ok {
			vv := v
			values[i] = &vv
		}
	}
	unit := defaultUnit
	field := ""
	pointId := 0
	if point != nil {
		field = point.Field
		pointId = point.Id
		if point.Unit != "" {
			unit = point.Unit
		}
	}
	return &sysin.PlcChartSeries{Name: name, Field: field, PointId: pointId, Type: "line", Data: values, Unit: unit}
}

func chartMeta(deviceId int, r plcChartRange, pointCount int) *sysin.PlcChartMeta {
	return &sysin.PlcChartMeta{DeviceId: deviceId, StartTime: r.Start, EndTime: r.End, Interval: r.Interval, PointCount: pointCount, MaxPoints: plcChartMaxPoints}
}

func appendChartSeries(series []*sysin.PlcChartSeries, item *sysin.PlcChartSeries) []*sysin.PlcChartSeries {
	if item == nil || item.PointId == 0 {
		return series
	}
	return append(series, item)
}

func chartLegend(series []*sysin.PlcChartSeries) []string {
	legend := make([]string, 0, len(series))
	for _, item := range series {
		legend = append(legend, item.Name)
	}
	return legend
}

// resolveField 未显式传 field 时，按当前设备已有点位自动匹配真实字段名（取第一个）。
func (s *sPlcChart) resolveField(ctx context.Context, deviceId int, candidates, excludes []string) string {
	matched := s.resolveFields(ctx, deviceId, candidates, excludes, true)
	if len(matched) == 0 {
		return ""
	}
	return matched[0]
}

// resolveFields 返回当前设备所有命中候选关键字的字段；firstOnly=true 时只返回第一个命中的候选组。
func (s *sPlcChart) resolveFields(ctx context.Context, deviceId int, candidates, excludes []string, firstOnly bool) []string {
	var rows []struct {
		Field string `orm:"field"`
	}
	_ = dao.PlcPoint.Ctx(ctx).
		Fields(dao.PlcPoint.Columns().Field).
		Where(dao.PlcPoint.Columns().DeviceId, deviceId).
		Where(dao.PlcPoint.Columns().Status, 1).
		OrderAsc(dao.PlcPoint.Columns().Sort).
		OrderAsc(dao.PlcPoint.Columns().Id).
		Scan(&rows)

	matched := make([]string, 0)
	seen := map[string]struct{}{}
	for _, candidate := range candidates {
		candidate = strings.ToLower(candidate)
		hit := false
		for _, row := range rows {
			field := row.Field
			lower := strings.ToLower(field)
			if !strings.Contains(lower, candidate) {
				continue
			}
			matchedExclude := false
			for _, exclude := range excludes {
				if strings.Contains(lower, strings.ToLower(exclude)) {
					matchedExclude = true
					break
				}
			}
			if matchedExclude {
				continue
			}
			if _, ok := seen[field]; ok {
				continue
			}
			seen[field] = struct{}{}
			matched = append(matched, field)
			hit = true
		}
		if firstOnly && hit {
			break
		}
	}
	return matched
}

// findPointByField 根据 field 找点位。
func (s *sPlcChart) findPointByField(ctx context.Context, deviceId int, field string) *plcChartPoint {
	if field == "" {
		return nil
	}
	p, err := s.findPoint(ctx, deviceId, field)
	if err != nil || p == nil {
		return nil
	}
	return p
}

// Temperature 单台设备的温度折线（设备温度 / 回油温度 / 油箱温度）。
func (s *sPlcChart) Temperature(ctx context.Context, in *sysin.PlcChartTemperatureInp) (res *sysin.PlcChartModel, err error) {
	r := defaultChartRange(ctx, in.DeviceId, in.StartTime, in.EndTime)

	excludeCommon := []string{"alarm", "al_", "sensor", "high", "low", "stop", "报警", "传感器", "高", "低", "停"}

	fieldOilBack := in.FieldOilBack
	if fieldOilBack == "" {
		fieldOilBack = s.resolveField(ctx, in.DeviceId,
			[]string{"lube_return_temp", "return_oil_temp", "return_temp", "回油温度", "回油温"},
			excludeCommon)
	}
	fieldOilTank := in.FieldOilTank
	if fieldOilTank == "" {
		fieldOilTank = s.resolveField(ctx, in.DeviceId,
			[]string{"lube_tank_temp", "oil_tank_temp", "tank_temp", "油箱温度", "油箱温"},
			excludeCommon)
	}
	// 设备温度（绕组/轴承/定子）：优先 1#WDG1，再退化到 WDG，再到 BEAR / STATOR / FRONT_BEARING；中文兼容“电机绕组/电机轴承/定子”。
	fieldTemp := in.FieldTemp
	if fieldTemp == "" {
		fieldTemp = s.resolveField(ctx, in.DeviceId,
			[]string{"1#mot_wdg1_temp", "mot_wdg1_temp", "stator_temp_a", "front_bearing_temp", "mot_wdg_temp", "mot_bear1_temp", "bearing_temp", "stator_temp", "device_temp", "电机绕组温度", "绕组温度", "定子温度", "电机轴承温度", "轴承温度", "电机温度", "设备温度"},
			excludeCommon)
	}

	oilBackPoint := s.findPointByField(ctx, in.DeviceId, fieldOilBack)
	oilTankPoint := s.findPointByField(ctx, in.DeviceId, fieldOilTank)
	tempPoint := s.findPointByField(ctx, in.DeviceId, fieldTemp)

	dataOilBack := map[string]float64{}
	if oilBackPoint != nil {
		if dataOilBack, err = s.queryPointData(ctx, oilBackPoint.Id, r); err != nil {
			return nil, err
		}
	}
	dataOilTank := map[string]float64{}
	if oilTankPoint != nil {
		if dataOilTank, err = s.queryPointData(ctx, oilTankPoint.Id, r); err != nil {
			return nil, err
		}
	}
	dataTemp := map[string]float64{}
	if tempPoint != nil {
		if dataTemp, err = s.queryPointData(ctx, tempPoint.Id, r); err != nil {
			return nil, err
		}
	}

	xAxis := collectXAxis(dataOilBack, dataOilTank, dataTemp)
	series := []*sysin.PlcChartSeries{}
	series = appendChartSeries(series, buildSeries("设备温度", tempPoint, "℃", xAxis, dataTemp))
	series = appendChartSeries(series, buildSeries("回油温度", oilBackPoint, "℃", xAxis, dataOilBack))
	series = appendChartSeries(series, buildSeries("油箱温度", oilTankPoint, "℃", xAxis, dataOilTank))

	return &sysin.PlcChartModel{Meta: chartMeta(in.DeviceId, r, len(xAxis)), Legend: chartLegend(series), XAxis: xAxis, Series: series}, nil
}

// Current 单台设备的电流折线；设备含多路电流时全部返回（如 1#/2# CRUSHER_AMPERE_ZSJ）。
func (s *sPlcChart) Current(ctx context.Context, in *sysin.PlcChartCurrentInp) (res *sysin.PlcChartModel, err error) {
	r := defaultChartRange(ctx, in.DeviceId, in.StartTime, in.EndTime)

	excludeCommon := []string{"alarm", "al_", "报警"}

	var fields []string
	if in.FieldCurrent != "" {
		fields = []string{in.FieldCurrent}
	} else {
		fields = s.resolveFields(ctx, in.DeviceId,
			[]string{"crusher_ampere", "main_current", "ampere", "current", "破碎机电流", "主电流", "电机电流", "电流"},
			excludeCommon, true)
	}

	allData := make([]map[string]float64, 0, len(fields))
	points := make([]*plcChartPoint, 0, len(fields))
	for _, f := range fields {
		p := s.findPointByField(ctx, in.DeviceId, f)
		data := map[string]float64{}
		if p != nil {
			if data, err = s.queryPointData(ctx, p.Id, r); err != nil {
				return nil, err
			}
		}
		points = append(points, p)
		allData = append(allData, data)
	}
	xAxis := collectXAxis(allData...)

	series := []*sysin.PlcChartSeries{}
	for i, p := range points {
		name := "电流"
		if p != nil && p.Field != "" {
			// 多路电流时用 field 作为序列名，单路时仍叫"电流"
			if len(points) > 1 {
				name = p.Field
			}
		}
		series = appendChartSeries(series, buildSeries(name, p, "A", xAxis, allData[i]))
	}
	return &sysin.PlcChartModel{Meta: chartMeta(in.DeviceId, r, len(xAxis)), Legend: chartLegend(series), XAxis: xAxis, Series: series}, nil
}
