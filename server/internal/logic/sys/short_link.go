// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
// @AutoGenerate Version 2.18.6
package sys

import (
	"context"
	"fmt"
	"net/url"
	apishortlink "hotgo/api/api/shortlink"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/library/hgorm/hook"
	"hotgo/internal/library/location"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
	"hotgo/utility/convert"
	"hotgo/utility/excel"
	"hotgo/utility/useragent"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

type sSysShortLink struct{}

func NewSysShortLink() *sSysShortLink {
	return &sSysShortLink{}
}

func init() {
	service.RegisterSysShortLink(NewSysShortLink())
}

// Model 短链接ORM模型
func (s *sSysShortLink) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.ShortLink.Ctx(ctx), option...)
}

// List 获取短链接列表
func (s *sSysShortLink) List(ctx context.Context, in *sysin.ShortLinkListInp) (list []*sysin.ShortLinkListModel, totalCount int, err error) {
	mod := s.Model(ctx)

	// 字段过滤
	mod = mod.Fields(sysin.ShortLinkListModel{})

	// 查询主键
	if in.Id > 0 {
		mod = mod.Where(dao.ShortLink.Columns().Id, in.Id)
	}

	// 查询状态：1正常 2禁用
	if in.Status > 0 {
		mod = mod.Where(dao.ShortLink.Columns().Status, in.Status)
	}

	// 查询创建时间
	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween(dao.ShortLink.Columns().CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	// 分页
	mod = mod.Page(in.Page, in.PerPage)

	// 排序
	mod = mod.OrderDesc(dao.ShortLink.Columns().Id)

	// 操作人摘要信息
	mod = mod.Hook(hook.MemberSummary)

	// 查询数据
	if err = mod.ScanAndCount(&list, &totalCount, false); err != nil {
		err = gerror.Wrap(err, "获取短链接列表失败，请稍后重试！")
		return
	}
	return
}

// Export 导出短链接
func (s *sSysShortLink) Export(ctx context.Context, in *sysin.ShortLinkListInp) (err error) {
	list, totalCount, err := s.List(ctx, in)
	if err != nil {
		return
	}

	// 字段的排序是依据tags的字段顺序，如果你不想使用默认的排序方式，可以直接定义 tags = []string{"字段名称", "字段名称2", ...}
	tags, err := convert.GetEntityDescTags(sysin.ShortLinkExportModel{})
	if err != nil {
		return
	}

	var (
		fileName  = "导出短链接-" + gctx.CtxId(ctx)
		sheetName = fmt.Sprintf("索引条件共%v行,共%v页,当前导出是第%v页,本页共%v行", totalCount, form.CalPageCount(totalCount, in.PerPage), in.Page, len(list))
		exports   []sysin.ShortLinkExportModel
	)

	if err = gconv.Scan(list, &exports); err != nil {
		return
	}

	err = excel.ExportByStructs(ctx, tags, exports, fileName, sheetName)
	return
}

// Create 前台创建短链接（自动生成短码）
func (s *sSysShortLink) Create(ctx context.Context, in *apishortlink.CreateReq) (res *apishortlink.CreateRes, err error) {
	// 生成唯一短码
	code, err := s.genUniqueCode(ctx)
	if err != nil {
		return
	}

	// 处理过期时间
	var expireAt *gtime.Time
	if in.ExpireAt != "" {
		expireAt = gtime.NewFromStr(in.ExpireAt)
	}

	// 写入数据库
	result, err := s.Model(ctx, &handler.Option{FilterAuth: false}).Data(g.Map{
		dao.ShortLink.Columns().Code:        code,
		dao.ShortLink.Columns().OriginalUrl: in.OriginalUrl,
		dao.ShortLink.Columns().Title:       in.Title,
		dao.ShortLink.Columns().ExpireAt:    expireAt,
		dao.ShortLink.Columns().Status:      1,
	}).Insert()
	if err != nil {
		err = gerror.Wrap(err, "创建短链接失败，请稍后重试！")
		return
	}

	id, _ := result.LastInsertId()
	domain := g.Cfg().MustGet(ctx, "shortlink.domain", "http://localhost:8002").String()
	res = &apishortlink.CreateRes{
		Id:       id,
		Code:     code,
		ShortUrl: domain + "/r/" + code,
	}
	return
}

// genUniqueCode 生成唯一短码（6位 base62，冲突重试5次）
func (s *sSysShortLink) genUniqueCode(ctx context.Context) (code string, err error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < 5; i++ {
		b := make([]byte, 6)
		for j := range b {
			b[j] = charset[grand.N(0, len(charset)-1)]
		}
		code = string(b)
		count, _ := s.Model(ctx).Where(dao.ShortLink.Columns().Code, code).Count()
		if count == 0 {
			return
		}
	}
	err = gerror.New("短码生成失败，请重试")
	return
}

// Redirect 前台获取原始链接并异步累加点击量
func (s *sSysShortLink) Redirect(ctx context.Context, code string) (res *apishortlink.RedirectRes, err error) {
	var record *entity.ShortLink
	if err = s.Model(ctx).Where(dao.ShortLink.Columns().Code, code).Scan(&record); err != nil {
		err = gerror.Wrap(err, "查询失败")
		return
	}
	if record == nil {
		err = gerror.New("短链接不存在")
		return
	}
	if record.Status != 1 {
		err = gerror.New("短链接已禁用")
		return
	}
	if record.ExpireAt != nil && record.ExpireAt.Before(gtime.Now()) {
		err = gerror.New("短链接已过期")
		return
	}

	// 提取请求信息（在 goroutine 外，ctx 还有效时读取）
	var (
		clientIp  string
		referer   string
		ua        string
	)
	if r := g.RequestFromCtx(ctx); r != nil {
		clientIp = location.GetClientIp(r)
		referer  = r.GetHeader("Referer")
		ua       = r.GetHeader("User-Agent")
	}

	// 异步写日志 + 累加点击量，不阻塞响应
	go func() {
		bgCtx := gctx.New()

		// 解析 referer host
		refererHost := ""
		if referer != "" {
			if u, e := url.Parse(referer); e == nil {
				refererHost = u.Host
			}
		}

		// 解析 UA
		browser := useragent.GetBrowser(ua)
		os      := useragent.GetOs(ua)

		// 解析 IP 归属地
		region  := ""
		country := ""
		if loc, e := location.GetLocation(bgCtx, clientIp); e == nil && loc != nil {
			country = loc.Country
			if loc.Province != "" || loc.City != "" {
				region = loc.Province + " " + loc.City
			}
		}

		// 写点击日志
		_, _ = dao.ShortLinkLog.Ctx(bgCtx).Data(g.Map{
			dao.ShortLinkLog.Columns().LinkId:      record.Id,
			dao.ShortLinkLog.Columns().Code:        record.Code,
			dao.ShortLinkLog.Columns().Ip:          clientIp,
			dao.ShortLinkLog.Columns().Region:      region,
			dao.ShortLinkLog.Columns().Country:     country,
			dao.ShortLinkLog.Columns().Referer:     referer,
			dao.ShortLinkLog.Columns().RefererHost: refererHost,
			dao.ShortLinkLog.Columns().UserAgent:   ua,
			dao.ShortLinkLog.Columns().Browser:     browser,
			dao.ShortLinkLog.Columns().Os:          os,
			dao.ShortLinkLog.Columns().CreatedAt:   gtime.Now(),
		}).Insert()

		// 累加点击量
		_, _ = s.Model(bgCtx).WherePri(record.Id).Increment(dao.ShortLink.Columns().TotalClicks, 1)
		_, _ = s.Model(bgCtx).WherePri(record.Id).Increment(dao.ShortLink.Columns().TodayClicks, 1)
		_, _ = s.Model(bgCtx).WherePri(record.Id).Increment(dao.ShortLink.Columns().WeeklyClicks, 1)
		_, _ = s.Model(bgCtx).WherePri(record.Id).Increment(dao.ShortLink.Columns().MonthlyClicks, 1)
	}()

	res = &apishortlink.RedirectRes{OriginalUrl: record.OriginalUrl}
	return
}

// Stats 前台获取短链接统计信息
func (s *sSysShortLink) Stats(ctx context.Context, code string) (res *apishortlink.StatsRes, err error) {
	var record *entity.ShortLink
	if err = s.Model(ctx).Where(dao.ShortLink.Columns().Code, code).Scan(&record); err != nil {
		err = gerror.Wrap(err, "查询失败")
		return
	}
	if record == nil {
		err = gerror.New("短链接不存在")
		return
	}

	domain := g.Cfg().MustGet(ctx, "shortlink.domain", "http://localhost:8002").String()
	res = &apishortlink.StatsRes{
		Id:              record.Id,
		Code:            record.Code,
		ShortUrl:        domain + "/r/" + record.Code,
		Title:           record.Title,
		OriginalUrl:     record.OriginalUrl,
		TotalClicks:     record.TotalClicks,
		TodayClicks:     record.TodayClicks,
		YesterdayClicks: record.YesterdayClicks,
		WeeklyClicks:    record.WeeklyClicks,
		MonthlyClicks:   record.MonthlyClicks,
	}
	if record.ExpireAt != nil {
		res.ExpireAt = record.ExpireAt.String()
	}
	if record.CreatedAt != nil {
		res.CreatedAt = record.CreatedAt.String()
	}

	logBase := dao.ShortLinkLog.Ctx(ctx).Where(dao.ShortLinkLog.Columns().LinkId, record.Id)

	// 近 14 天访问趋势（避免用 date 保留字作别名）
	type trendRow struct {
		ClickDate string `orm:"click_date"`
		Clicks    int    `orm:"clicks"`
	}
	var trendRows []trendRow
	_ = logBase.
		Fields("DATE(created_at) AS click_date, COUNT(*) AS clicks").
		Where("created_at >= ?", gtime.Now().AddDate(0, 0, -13).Format("Y-m-d 00:00:00")).
		Group("DATE(created_at)").
		Order("click_date ASC").
		Scan(&trendRows)
	res.Trend = make([]apishortlink.StatsTrendItem, 0, len(trendRows))
	for _, r := range trendRows {
		res.Trend = append(res.Trend, apishortlink.StatsTrendItem{Date: r.ClickDate, Clicks: r.Clicks})
	}

	// 来源域名 Top10
	type refRow struct {
		Source string `orm:"source"`
		Count  int    `orm:"cnt"`
	}
	var refRows []refRow
	_ = dao.ShortLinkLog.Ctx(ctx).
		Where(dao.ShortLinkLog.Columns().LinkId, record.Id).
		Fields("referer_host AS source, COUNT(*) AS cnt").
		Group("referer_host").
		Order("cnt DESC").
		Limit(10).
		Scan(&refRows)
	res.Referrers = make([]apishortlink.StatsReferrerItem, 0, len(refRows))
	for _, r := range refRows {
		res.Referrers = append(res.Referrers, apishortlink.StatsReferrerItem{Source: r.Source, Count: r.Count})
	}

	// 访问地区 Top10
	type regRow struct {
		Region string `orm:"region"`
		Count  int    `orm:"cnt"`
	}
	var regRows []regRow
	_ = dao.ShortLinkLog.Ctx(ctx).
		Where(dao.ShortLinkLog.Columns().LinkId, record.Id).
		WhereNot(dao.ShortLinkLog.Columns().Region, "").
		Fields("region, COUNT(*) AS cnt").
		Group("region").
		Order("cnt DESC").
		Limit(10).
		Scan(&regRows)
	res.Regions = make([]apishortlink.StatsRegionItem, 0, len(regRows))
	for _, r := range regRows {
		res.Regions = append(res.Regions, apishortlink.StatsRegionItem{Region: r.Region, Count: r.Count})
	}

	return
}

// ListPub 前台分页查询短链接列表
func (s *sSysShortLink) ListPub(ctx context.Context, in *apishortlink.ListReq) (res *apishortlink.ListRes, err error) {
	mod := s.Model(ctx).Where(dao.ShortLink.Columns().Status, 1)

	if in.Keyword != "" {
		kw := "%" + in.Keyword + "%"
		mod = mod.WhereOrLike(dao.ShortLink.Columns().Title, kw).
			WhereOrLike(dao.ShortLink.Columns().OriginalUrl, kw)
	}

	var (
		list  []*entity.ShortLink
		total int
	)
	if err = mod.Clone().OrderDesc(dao.ShortLink.Columns().Id).
		Page(in.Page, in.Limit).ScanAndCount(&list, &total, true); err != nil {
		err = gerror.Wrap(err, "获取短链接列表失败，请稍后重试！")
		return
	}

	domain := g.Cfg().MustGet(ctx, "shortlink.domain", "http://localhost:8002").String()
	items := make([]*apishortlink.ListItem, 0, len(list))
	for _, r := range list {
		item := &apishortlink.ListItem{
			Id:          r.Id,
			Code:        r.Code,
			ShortUrl:    domain + "/r/" + r.Code,
			OriginalUrl: r.OriginalUrl,
			Title:       r.Title,
			TotalClicks: r.TotalClicks,
			TodayClicks: r.TodayClicks,
		}
		if r.ExpireAt != nil {
			item.ExpireAt = r.ExpireAt.String()
		}
		if r.CreatedAt != nil {
			item.CreatedAt = r.CreatedAt.String()
		}
		items = append(items, item)
	}

	res = &apishortlink.ListRes{
		List:  items,
		Total: total,
		Page:  in.Page,
		Limit: in.Limit,
	}
	return
}

// Edit 修改/新增短链接
func (s *sSysShortLink) Edit(ctx context.Context, in *sysin.ShortLinkEditInp) (err error) {
	// 验证'Code'唯一
	if err = hgorm.IsUnique(ctx, &dao.ShortLink, g.Map{dao.ShortLink.Columns().Code: in.Code}, "短码已存在", in.Id); err != nil {
		return
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {

		// 修改
		if in.Id > 0 {
			in.UpdatedBy = contexts.GetUserId(ctx)
			if _, err = s.Model(ctx).
				Fields(sysin.ShortLinkUpdateFields{}).
				WherePri(in.Id).Data(in).Update(); err != nil {
				err = gerror.Wrap(err, "修改短链接失败，请稍后重试！")
			}
			return
		}

		// 新增
		in.CreatedBy = contexts.GetUserId(ctx)
		if _, err = s.Model(ctx, &handler.Option{FilterAuth: false}).
			Fields(sysin.ShortLinkInsertFields{}).
			Data(in).OmitEmptyData().Insert(); err != nil {
			err = gerror.Wrap(err, "新增短链接失败，请稍后重试！")
		}
		return
	})
}

// Delete 删除短链接
func (s *sSysShortLink) Delete(ctx context.Context, in *sysin.ShortLinkDeleteInp) (err error) {

	if _, err = s.Model(ctx).WherePri(in.Id).Unscoped().Delete(); err != nil {
		err = gerror.Wrap(err, "删除短链接失败，请稍后重试！")
		return
	}
	return
}

// View 获取短链接指定信息
func (s *sSysShortLink) View(ctx context.Context, in *sysin.ShortLinkViewInp) (res *sysin.ShortLinkViewModel, err error) {
	if err = s.Model(ctx).WherePri(in.Id).Hook(hook.MemberSummary).Scan(&res); err != nil {
		err = gerror.Wrap(err, "获取短链接信息，请稍后重试！")
		return
	}
	return
}

// Status 更新短链接状态
func (s *sSysShortLink) Status(ctx context.Context, in *sysin.ShortLinkStatusInp) (err error) {
	if _, err = s.Model(ctx).WherePri(in.Id).Data(g.Map{
		dao.ShortLink.Columns().Status:    in.Status,
		dao.ShortLink.Columns().UpdatedBy: contexts.GetUserId(ctx),
	}).Update(); err != nil {
		err = gerror.Wrap(err, "更新短链接状态失败，请稍后重试！")
		return
	}
	return
}