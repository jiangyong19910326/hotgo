// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ShortLinkLogDao is the data access object for the table hg_short_link_log.
type ShortLinkLogDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  ShortLinkLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// ShortLinkLogColumns defines and stores column names for the table hg_short_link_log.
type ShortLinkLogColumns struct {
	Id          string // 主键
	LinkId      string // 短链ID
	Code        string // 短码（冗余，提高查询效率）
	Ip          string // 访客 IP
	Region      string // 访问地区（省市）
	Country     string // 访问国家
	Referer     string // 来源完整 URL
	RefererHost string // 来源域名
	UserAgent   string // User-Agent
	Device      string // 设备类型：desktop / mobile / tablet
	Browser     string // 浏览器
	Os          string // 操作系统
	CreatedAt   string // 点击时间
}

// shortLinkLogColumns holds the columns for the table hg_short_link_log.
var shortLinkLogColumns = ShortLinkLogColumns{
	Id:          "id",
	LinkId:      "link_id",
	Code:        "code",
	Ip:          "ip",
	Region:      "region",
	Country:     "country",
	Referer:     "referer",
	RefererHost: "referer_host",
	UserAgent:   "user_agent",
	Device:      "device",
	Browser:     "browser",
	Os:          "os",
	CreatedAt:   "created_at",
}

// NewShortLinkLogDao creates and returns a new DAO object for table data access.
func NewShortLinkLogDao(handlers ...gdb.ModelHandler) *ShortLinkLogDao {
	return &ShortLinkLogDao{
		group:    "default",
		table:    "hg_short_link_log",
		columns:  shortLinkLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ShortLinkLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ShortLinkLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ShortLinkLogDao) Columns() ShortLinkLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ShortLinkLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ShortLinkLogDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *ShortLinkLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
