package plc

import (
	v1 "hotgo/api/api/plc"
	"hotgo/internal/model/input/sysin"
)

func toV1Chart(src *sysin.PlcChartModel) *v1.ChartModel {
	if src == nil {
		return &v1.ChartModel{XAxis: []string{}, Series: []*v1.ChartSeries{}}
	}
	return &v1.ChartModel{
		Meta:   toV1Meta(src.Meta),
		Legend: src.Legend,
		XAxis:  src.XAxis,
		Series: toV1Series(src.Series),
	}
}

func toV1Meta(src *sysin.PlcChartMeta) *v1.ChartMeta {
	if src == nil {
		return nil
	}
	return &v1.ChartMeta{
		DeviceId:   src.DeviceId,
		StartTime:  src.StartTime,
		EndTime:    src.EndTime,
		Interval:   src.Interval,
		PointCount: src.PointCount,
		MaxPoints:  src.MaxPoints,
	}
}

// toV1Series 将 sysin.PlcChartSeries 转换为前台响应的 v1.ChartSeries
func toV1Series(src []*sysin.PlcChartSeries) []*v1.ChartSeries {
	out := make([]*v1.ChartSeries, 0, len(src))
	for _, s := range src {
		out = append(out, &v1.ChartSeries{
			Name:    s.Name,
			Field:   s.Field,
			PointId: s.PointId,
			Type:    s.Type,
			Data:    s.Data,
			Unit:    s.Unit,
		})
	}
	return out
}
