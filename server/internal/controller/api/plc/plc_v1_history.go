package plc

import (
	"context"

	v1 "hotgo/api/api/plc"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

const (
	historyDefaultLimit = 500
	historyMaxLimit     = 5000
)

func (c *ControllerV1) History(ctx context.Context, req *v1.HistoryReq) (res *v1.HistoryRes, err error) {
	limit := req.Limit
	if limit <= 0 {
		limit = historyDefaultLimit
	}
	if limit > historyMaxLimit {
		limit = historyMaxLimit
	}

	end := gtime.Now()
	if req.EndTime != "" {
		end = gtime.NewFromStr(req.EndTime)
		if end == nil {
			end = gtime.Now()
		}
	}
	start := end.Add(-1 * 60 * 60 * 1000 * 1000 * 1000) // -1h
	if req.StartTime != "" {
		s := gtime.NewFromStr(req.StartTime)
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
			PointId: point.Id,
			Field:   point.Field,
			Name:    point.Name,
			Unit:    point.Unit,
		},
		List: []*v1.HistoryItem{},
	}
	if point.Id == 0 {
		return res, nil
	}

	// 查记录, 升序
	var rows []*entity.PlcRecord
	err = dao.PlcRecord.Ctx(ctx).
		Where(dao.PlcRecord.Columns().PointId, req.PointId).
		WhereGTE(dao.PlcRecord.Columns().CollectedAt, start).
		WhereLTE(dao.PlcRecord.Columns().CollectedAt, end).
		OrderAsc(dao.PlcRecord.Columns().CollectedAt).
		Limit(limit).
		Scan(&rows)
	if err != nil {
		return nil, err
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
