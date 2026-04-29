// Package crons PLC 历史记录清理
package crons

import (
	"context"

	"hotgo/internal/dao"
	"hotgo/internal/library/cron"

	"github.com/gogf/gf/v2/os/gtime"
)

const plcRecordRetainDays = 30

func init() {
	cron.Register(PlcCleanup)
}

// PlcCleanup 清理 30 天前的 PLC 采集历史记录
// 在系统管理 → 定时任务 中配置 Pattern, 例如 "0 0 3 * * *" (每天 03:00)
var PlcCleanup = &cPlcCleanup{name: "plc_cleanup"}

type cPlcCleanup struct {
	name string
}

func (c *cPlcCleanup) GetName() string { return c.name }

func (c *cPlcCleanup) Execute(ctx context.Context, parser *cron.Parser) (err error) {
	cutoff := gtime.Now().AddDate(0, 0, -plcRecordRetainDays)
	res, err := dao.PlcRecord.Ctx(ctx).
		WhereLT(dao.PlcRecord.Columns().CollectedAt, cutoff).
		Delete()
	if err != nil {
		parser.Logger.Warningf(ctx, "plc_cleanup err: %v", err)
		return nil
	}
	n, _ := res.RowsAffected()
	parser.Logger.Infof(ctx, "plc_cleanup deleted %d records older than %s", n, cutoff.Format("Y-m-d H:i:s"))
	return
}
