package plc

import (
	"context"

	v1 "hotgo/api/open/plc"
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

func (c *ControllerV1) History(ctx context.Context, req *v1.HistoryReq) (res *v1.HistoryRes, err error) {
	perPage := req.PerPage
	if perPage <= 0 {
		perPage = 200
	}
	if perPage > 500 {
		perPage = 500
	}

	list, totalCount, err := service.PlcHistory().List(ctx, &sysin.PlcHistoryInp{
		PageReq:   form.PageReq{Page: req.Page, PerPage: perPage},
		PointId:   req.PointId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	res = &v1.HistoryRes{
		TotalCount: totalCount,
		Page:       req.Page,
		PerPage:    perPage,
		List:       make([]*v1.HistoryItem, 0, len(list)),
	}
	for _, h := range list {
		item := &v1.HistoryItem{
			EngValue: h.EngValue,
			RawValue: h.RawValue,
		}
		if h.CollectedAt != nil {
			item.CollectedAt = h.CollectedAt.String()
		}
		res.List = append(res.List, item)
	}
	return
}
