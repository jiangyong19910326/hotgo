package plc

import (
	"context"

	v1 "hotgo/api/open/plc"
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
)

func (c *ControllerV1) AlarmList(ctx context.Context, req *v1.AlarmListReq) (res *v1.AlarmListRes, err error) {
	perPage := req.PerPage
	if perPage <= 0 {
		perPage = 20
	}

	list, totalCount, err := service.PlcAlarm().List(ctx, &sysin.PlcAlarmListInp{
		PageReq:    form.PageReq{Page: req.Page, PerPage: perPage},
		DeviceId:   req.DeviceId,
		IsResolved: req.IsResolved,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	res = &v1.AlarmListRes{
		TotalCount: totalCount,
		Page:       req.Page,
		PerPage:    perPage,
		List:       make([]*v1.AlarmItem, 0, len(list)),
	}
	for _, a := range list {
		item := &v1.AlarmItem{
			Id:        int(a.Id),
			DeviceId:  a.DeviceId,
			PointId:   a.PointId,
			PointName: a.PointName,
			EngValue:  a.EngValue,
			AlarmType: a.AlarmType,
			AlarmMin:  a.AlarmMin,
			AlarmMax:  a.AlarmMax,
			Unit:      a.Unit,
			IsResolved: a.IsResolved,
			Remark:    a.Remark,
		}
		if a.ResolvedAt != nil {
			item.ResolvedAt = a.ResolvedAt.String()
		}
		if a.TriggeredAt != nil {
			item.TriggeredAt = a.TriggeredAt.String()
		}
		res.List = append(res.List, item)
	}
	return
}
