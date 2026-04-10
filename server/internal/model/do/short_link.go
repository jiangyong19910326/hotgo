// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShortLink is the golang structure of table hg_short_link for DAO operations like Where/Data.
type ShortLink struct {
	g.Meta          `orm:"table:hg_short_link, do:true"`
	Id              any         // 主键
	Code            any         // 短码
	OriginalUrl     any         // 原始链接
	Title           any         // 标题
	TotalClicks     any         // 总点击量
	TodayClicks     any         // 今日点击（每日凌晨重置）
	YesterdayClicks any         // 昨日点击
	WeeklyClicks    any         // 本周点击
	MonthlyClicks   any         // 本月点击
	ExpireAt        *gtime.Time // 过期时间，NULL 表示永不过期
	Status          any         // 状态：1正常 2禁用
	CreatedBy       any         // 创建者 member_id
	UpdatedBy       any         // 更新者 member_id
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 修改时间
	DeletedAt       *gtime.Time // 删除时间（软删除）
}
