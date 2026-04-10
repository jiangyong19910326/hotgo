// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ShortLink is the golang structure for table short_link.
type ShortLink struct {
	Id              int64       `json:"id"              orm:"id"               description:"主键"`
	Code            string      `json:"code"            orm:"code"             description:"短码"`
	OriginalUrl     string      `json:"originalUrl"     orm:"original_url"     description:"原始链接"`
	Title           string      `json:"title"           orm:"title"            description:"标题"`
	TotalClicks     int64       `json:"totalClicks"     orm:"total_clicks"     description:"总点击量"`
	TodayClicks     int         `json:"todayClicks"     orm:"today_clicks"     description:"今日点击（每日凌晨重置）"`
	YesterdayClicks int         `json:"yesterdayClicks" orm:"yesterday_clicks" description:"昨日点击"`
	WeeklyClicks    int         `json:"weeklyClicks"    orm:"weekly_clicks"    description:"本周点击"`
	MonthlyClicks   int         `json:"monthlyClicks"   orm:"monthly_clicks"   description:"本月点击"`
	ExpireAt        *gtime.Time `json:"expireAt"        orm:"expire_at"        description:"过期时间，NULL 表示永不过期"`
	Status          int         `json:"status"          orm:"status"           description:"状态：1正常 2禁用"`
	CreatedBy       int64       `json:"createdBy"       orm:"created_by"       description:"创建者 member_id"`
	UpdatedBy       int64       `json:"updatedBy"       orm:"updated_by"       description:"更新者 member_id"`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:"修改时间"`
	DeletedAt       *gtime.Time `json:"deletedAt"       orm:"deleted_at"       description:"删除时间（软删除）"`
}
