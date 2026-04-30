// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FrontUserDao is the data access object for the table hg_front_user.
type FrontUserDao struct {
	table    string
	group    string
	columns  FrontUserColumns
	handlers []gdb.ModelHandler
}

// FrontUserColumns defines and stores column names for the table hg_front_user.
type FrontUserColumns struct {
	Id           string
	Username     string
	Nickname     string
	PasswordHash string
	Salt         string
	Avatar       string
	Mobile       string
	Email        string
	Remark       string
	Status       string
	LastLoginAt  string
	LastLoginIp  string
	CreatedBy    string
	UpdatedBy    string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}

var frontUserColumns = FrontUserColumns{
	Id:           "id",
	Username:     "username",
	Nickname:     "nickname",
	PasswordHash: "password_hash",
	Salt:         "salt",
	Avatar:       "avatar",
	Mobile:       "mobile",
	Email:        "email",
	Remark:       "remark",
	Status:       "status",
	LastLoginAt:  "last_login_at",
	LastLoginIp:  "last_login_ip",
	CreatedBy:    "created_by",
	UpdatedBy:    "updated_by",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

func NewFrontUserDao(handlers ...gdb.ModelHandler) *FrontUserDao {
	return &FrontUserDao{
		group:    "default",
		table:    "hg_front_user",
		columns:  frontUserColumns,
		handlers: handlers,
	}
}

func (dao *FrontUserDao) DB() gdb.DB                { return g.DB(dao.group) }
func (dao *FrontUserDao) Table() string             { return dao.table }
func (dao *FrontUserDao) Columns() FrontUserColumns { return dao.columns }
func (dao *FrontUserDao) Group() string             { return dao.group }

func (dao *FrontUserDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, h := range dao.handlers {
		model = h(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *FrontUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
