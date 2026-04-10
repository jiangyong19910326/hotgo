// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
// @AutoGenerate Version 2.18.6
package sysin

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/library/hgorm/hook"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"
	"hotgo/utility/validate"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShortLinkUpdateFields 修改短链接字段过滤
type ShortLinkUpdateFields struct {
	Code            string      `json:"code"            dc:"短码"`
	OriginalUrl     string      `json:"originalUrl"     dc:"原始链接"`
	Title           string      `json:"title"           dc:"标题"`
	TotalClicks     int64       `json:"totalClicks"     dc:"总点击量"`
	TodayClicks     int         `json:"todayClicks"     dc:"今日点击（每日凌晨重置）"`
	YesterdayClicks int         `json:"yesterdayClicks" dc:"昨日点击"`
	WeeklyClicks    int         `json:"weeklyClicks"    dc:"本周点击"`
	MonthlyClicks   int         `json:"monthlyClicks"   dc:"本月点击"`
	ExpireAt        *gtime.Time `json:"expireAt"        dc:"过期时间，NULL 表示永不过期"`
	Status          int         `json:"status"          dc:"状态：1正常 2禁用"`
	UpdatedBy       int64       `json:"updatedBy"       dc:"更新者 member_id"`
}

// ShortLinkInsertFields 新增短链接字段过滤
type ShortLinkInsertFields struct {
	Code            string      `json:"code"            dc:"短码"`
	OriginalUrl     string      `json:"originalUrl"     dc:"原始链接"`
	Title           string      `json:"title"           dc:"标题"`
	TotalClicks     int64       `json:"totalClicks"     dc:"总点击量"`
	TodayClicks     int         `json:"todayClicks"     dc:"今日点击（每日凌晨重置）"`
	YesterdayClicks int         `json:"yesterdayClicks" dc:"昨日点击"`
	WeeklyClicks    int         `json:"weeklyClicks"    dc:"本周点击"`
	MonthlyClicks   int         `json:"monthlyClicks"   dc:"本月点击"`
	ExpireAt        *gtime.Time `json:"expireAt"        dc:"过期时间，NULL 表示永不过期"`
	Status          int         `json:"status"          dc:"状态：1正常 2禁用"`
	CreatedBy       int64       `json:"createdBy"       dc:"创建者 member_id"`
}

// ShortLinkEditInp 修改/新增短链接
type ShortLinkEditInp struct {
	entity.ShortLink
}

func (in *ShortLinkEditInp) Filter(ctx context.Context) (err error) {
	// 验证短码
	if err := g.Validator().Rules("required").Data(in.Code).Messages("短码不能为空").Run(ctx); err != nil {
		return err.Current()
	}
	// 验证原始链接
	if err := g.Validator().Rules("required").Data(in.OriginalUrl).Messages("原始链接不能为空").Run(ctx); err != nil {
		return err.Current()
	}
	return
}

type ShortLinkEditModel struct{}

// ShortLinkDeleteInp 删除短链接
type ShortLinkDeleteInp struct {
	Id interface{} `json:"id" v:"required#主键不能为空" dc:"主键"`
}

func (in *ShortLinkDeleteInp) Filter(ctx context.Context) (err error) {
	return
}

type ShortLinkDeleteModel struct{}

// ShortLinkViewInp 获取指定短链接信息
type ShortLinkViewInp struct {
	Id int64 `json:"id" v:"required#主键不能为空" dc:"主键"`
}

func (in *ShortLinkViewInp) Filter(ctx context.Context) (err error) {
	return
}

type ShortLinkViewModel struct {
	entity.ShortLink
	CreatedBySumma *hook.MemberSumma `json:"createdBySumma" dc:"创建者 member_id摘要信息"`
	UpdatedBySumma *hook.MemberSumma `json:"updatedBySumma" dc:"更新者 member_id摘要信息"`
}

// ShortLinkListInp 获取短链接列表
type ShortLinkListInp struct {
	form.PageReq
	Id        int64         `json:"id"        dc:"主键"`
	Status    int           `json:"status"    dc:"状态：1正常 2禁用"`
	CreatedAt []*gtime.Time `json:"createdAt" dc:"创建时间"`
}

func (in *ShortLinkListInp) Filter(ctx context.Context) (err error) {
	return
}

type ShortLinkListModel struct {
	Id              int64             `json:"id"              dc:"主键"`
	Code            string            `json:"code"            dc:"短码"`
	Title           string            `json:"title"           dc:"标题"`
	TotalClicks     int64             `json:"totalClicks"     dc:"总点击量"`
	TodayClicks     int               `json:"todayClicks"     dc:"今日点击（每日凌晨重置）"`
	YesterdayClicks int               `json:"yesterdayClicks" dc:"昨日点击"`
	WeeklyClicks    int               `json:"weeklyClicks"    dc:"本周点击"`
	MonthlyClicks   int               `json:"monthlyClicks"   dc:"本月点击"`
	ExpireAt        *gtime.Time       `json:"expireAt"        dc:"过期时间，NULL 表示永不过期"`
	Status          int               `json:"status"          dc:"状态：1正常 2禁用"`
	CreatedBy       int64             `json:"createdBy"       dc:"创建者 member_id"`
	CreatedBySumma  *hook.MemberSumma `json:"createdBySumma"  dc:"创建者 member_id摘要信息"`
	UpdatedBy       int64             `json:"updatedBy"       dc:"更新者 member_id"`
	UpdatedBySumma  *hook.MemberSumma `json:"updatedBySumma"  dc:"更新者 member_id摘要信息"`
	CreatedAt       *gtime.Time       `json:"createdAt"       dc:"创建时间"`
	UpdatedAt       *gtime.Time       `json:"updatedAt"       dc:"修改时间"`
}

// ShortLinkExportModel 导出短链接
type ShortLinkExportModel struct {
	Id              int64       `json:"id"              dc:"主键"`
	Code            string      `json:"code"            dc:"短码"`
	Title           string      `json:"title"           dc:"标题"`
	TotalClicks     int64       `json:"totalClicks"     dc:"总点击量"`
	TodayClicks     int         `json:"todayClicks"     dc:"今日点击（每日凌晨重置）"`
	YesterdayClicks int         `json:"yesterdayClicks" dc:"昨日点击"`
	WeeklyClicks    int         `json:"weeklyClicks"    dc:"本周点击"`
	MonthlyClicks   int         `json:"monthlyClicks"   dc:"本月点击"`
	ExpireAt        *gtime.Time `json:"expireAt"        dc:"过期时间，NULL 表示永不过期"`
	Status          int         `json:"status"          dc:"状态：1正常 2禁用"`
	CreatedBy       int64       `json:"createdBy"       dc:"创建者 member_id"`
	UpdatedBy       int64       `json:"updatedBy"       dc:"更新者 member_id"`
	CreatedAt       *gtime.Time `json:"createdAt"       dc:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updatedAt"       dc:"修改时间"`
}

// ShortLinkStatusInp 更新短链接状态
type ShortLinkStatusInp struct {
	Id     int64 `json:"id" v:"required#主键不能为空" dc:"主键"`
	Status int   `json:"status" dc:"状态"`
}

func (in *ShortLinkStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		err = gerror.New("主键不能为空")
		return
	}

	if in.Status <= 0 {
		err = gerror.New("状态不能为空")
		return
	}

	if !validate.InSlice(consts.StatusSlice, in.Status) {
		err = gerror.New("状态不正确")
		return
	}
	return
}

type ShortLinkStatusModel struct{}