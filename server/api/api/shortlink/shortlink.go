package shortlink

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CreateReq 创建短链接
type CreateReq struct {
	g.Meta      `path:"/shortlink/create" method:"post" tags:"短链接" summary:"创建短链接"`
	OriginalUrl string `json:"original_url" v:"required|url#原始链接不能为空|链接格式不正确" dc:"原始链接"`
	Title       string `json:"title"        dc:"标题"`
	ExpireAt    string `json:"expire_at"    dc:"过期时间，空表示永不过期"`
}

type CreateRes struct {
	Id       int64  `json:"id"`
	Code     string `json:"code"`
	ShortUrl string `json:"short_url"`
}

// RedirectReq 获取原始链接（用于前台跳转）
type RedirectReq struct {
	g.Meta `path:"/shortlink/redirect/{code}" method:"get" tags:"短链接" summary:"获取短链接跳转地址"`
	Code   string `p:"code" v:"required" dc:"短码"`
}

type RedirectRes struct {
	OriginalUrl string `json:"original_url"`
}

// StatsReq 获取统计信息
type StatsReq struct {
	g.Meta `path:"/shortlink/stats/{code}" method:"get" tags:"短链接" summary:"获取短链接统计"`
	Code   string `p:"code" v:"required" dc:"短码"`
}

type StatsTrendItem struct {
	Date   string `json:"date"`
	Clicks int    `json:"clicks"`
}

type StatsReferrerItem struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

type StatsRegionItem struct {
	Region string `json:"region"`
	Count  int    `json:"count"`
}

type StatsRes struct {
	Id              int64               `json:"id"`
	Code            string              `json:"code"`
	ShortUrl        string              `json:"short_url"`
	Title           string              `json:"title"`
	OriginalUrl     string              `json:"original_url"`
	TotalClicks     int64               `json:"total_clicks"`
	TodayClicks     int                 `json:"today_clicks"`
	YesterdayClicks int                 `json:"yesterday_clicks"`
	WeeklyClicks    int                 `json:"weekly_clicks"`
	MonthlyClicks   int                 `json:"monthly_clicks"`
	ExpireAt        string              `json:"expire_at"`
	CreatedAt       string              `json:"created_at"`
	Trend           []StatsTrendItem    `json:"trend"`
	Referrers       []StatsReferrerItem `json:"referrers"`
	Regions         []StatsRegionItem   `json:"regions"`
}

// ListReq 分页列表
type ListReq struct {
	g.Meta  `path:"/shortlink/list" method:"get" tags:"短链接" summary:"短链接列表"`
	Page    int    `json:"page"    d:"1"`
	Limit   int    `json:"limit"   d:"15"`
	Keyword string `json:"keyword" dc:"关键词搜索"`
}

type ListItem struct {
	Id          int64  `json:"id"`
	Code        string `json:"code"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
	Title       string `json:"title"`
	TotalClicks int64  `json:"total_clicks"`
	TodayClicks int    `json:"today_clicks"`
	ExpireAt    string `json:"expire_at"`
	CreatedAt   string `json:"created_at"`
}

type ListRes struct {
	List  []*ListItem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}
