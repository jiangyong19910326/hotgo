package plc

import (
	"context"
	"strings"
	"time"

	v1 "hotgo/api/api/plc"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	historyDefaultLimit = 100
	historyMaxLimit     = 5000
)

func normalizeHistoryTime(value string, isStart bool) string {
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

func (c *ControllerV1) History(ctx context.Context, req *v1.HistoryReq) (res *v1.HistoryRes, err error) {
	if err = ensureDeviceAccess(ctx, req.DeviceId); err != nil {
		return nil, err
	}
	if err = ensurePointAccess(ctx, req.PointId); err != nil {
		return nil, err
	}
	limit := req.Limit
	if limit <= 0 {
		limit = historyDefaultLimit
	}
	if limit > historyMaxLimit {
		limit = historyMaxLimit
	}

	var start, end *gtime.Time
	if req.EndTime != "" {
		end = gtime.NewFromStr(normalizeHistoryTime(req.EndTime, false))
	}
	if req.StartTime != "" {
		s := gtime.NewFromStr(normalizeHistoryTime(req.StartTime, true))
		if s != nil {
			start = s
		}
	}

	// 取点位元信息
	var point entity.PlcPoint
	if perr := dao.PlcPoint.Ctx(ctx).Where(dao.PlcPoint.Columns().Id, req.PointId).Scan(&point); perr != nil {
		return nil, perr
	}
	res = &v1.HistoryRes{
		Point: &v1.HistoryPointInfo{
			PointId:  point.Id,
			DeviceId: point.DeviceId,
			Field:    point.Field,
			Name:     point.Name,
			Unit:     point.Unit,
		},
		List: []*v1.HistoryItem{},
	}
	if point.Id == 0 {
		return res, nil
	}
	if point.DeviceId != req.DeviceId {
		return nil, gerror.New("点位不属于当前设备")
	}

	// 查记录：未传时间时取最新 limit 条；传时间时按范围过滤。
	var rows []*entity.PlcRecord
	mod := dao.PlcRecord.Ctx(ctx).Where(dao.PlcRecord.Columns().PointId, req.PointId)
	if start != nil {
		mod = mod.WhereGTE(dao.PlcRecord.Columns().CollectedAt, start)
	}
	if end != nil {
		mod = mod.WhereLTE(dao.PlcRecord.Columns().CollectedAt, end)
	}
	if start == nil && end == nil {
		mod = mod.OrderDesc(dao.PlcRecord.Columns().CollectedAt)
	} else {
		mod = mod.OrderAsc(dao.PlcRecord.Columns().CollectedAt)
	}
	err = mod.Limit(limit).Scan(&rows)
	if err != nil {
		return nil, err
	}
	if start == nil && end == nil {
		for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
			rows[i], rows[j] = rows[j], rows[i]
		}
	}

	res.List = make([]*v1.HistoryItem, 0, len(rows))
	for _, r := range rows {
		t := ""
		if r.CollectedAt != nil {
			t = r.CollectedAt.Format("Y-m-d H:i:s")
		}
		res.List = append(res.List, &v1.HistoryItem{Time: t, EngValue: r.EngValue})
	}
	return res, nil
}
