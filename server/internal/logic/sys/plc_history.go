// Package sys PLC 历史记录
package sys

import (
	"context"

	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

const plcRecordRetentionMonths = -3

type sPlcHistory struct{}

func NewPlcHistory() *sPlcHistory { return &sPlcHistory{} }

func init() {
	service.RegisterPlcHistory(NewPlcHistory())
}

func (s *sPlcHistory) List(ctx context.Context, in *sysin.PlcHistoryInp) (list []*sysin.PlcHistoryModel, totalCount int, err error) {
	mod := dao.PlcRecord.Ctx(ctx).Where(dao.PlcRecord.Columns().PointId, in.PointId)
	if in.StartTime != "" {
		mod = mod.WhereGTE(dao.PlcRecord.Columns().CollectedAt, in.StartTime)
	}
	if in.EndTime != "" {
		mod = mod.WhereLTE(dao.PlcRecord.Columns().CollectedAt, in.EndTime)
	}

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	err = mod.Page(in.Page, in.PerPage).
		OrderDesc(dao.PlcRecord.Columns().CollectedAt).
		Fields(
			dao.PlcRecord.Columns().CollectedAt,
			dao.PlcRecord.Columns().EngValue,
			dao.PlcRecord.Columns().RawValue,
		).
		Scan(&list)
	return
}

// BatchInsert 批量写入采集记录
func (s *sPlcHistory) BatchInsert(ctx context.Context, records []*entity.PlcRecord) error {
	if len(records) == 0 {
		return nil
	}
	now := gtime.Now()
	for _, r := range records {
		if r.CollectedAt == nil {
			r.CollectedAt = now
		}
	}
	_, err := dao.PlcRecord.Ctx(ctx).Data(records).Insert()
	if err != nil {
		return err
	}
	return s.CleanupExpired(ctx)
}

// CleanupExpired 清理三个月前的历史记录；仅清理三个月内仍有新数据的点位。
func (s *sPlcHistory) CleanupExpired(ctx context.Context) error {
	cutoff := gtime.Now().AddDate(0, plcRecordRetentionMonths, 0).Format("Y-m-d H:i:s")
	sql := `
DELETE r FROM hg_plc_record r
INNER JOIN (
  SELECT point_id
  FROM hg_plc_record
  GROUP BY point_id
  HAVING MAX(collected_at) >= ?
) active_points ON active_points.point_id = r.point_id
WHERE r.collected_at < ?`
	_, err := dao.PlcRecord.DB().Exec(ctx, sql, cutoff, cutoff)
	return err
}
