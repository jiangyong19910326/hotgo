// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShortLinkLog is the golang structure of table hg_short_link_log for DAO operations like Where/Data.
type ShortLinkLog struct {
	g.Meta      `orm:"table:hg_short_link_log, do:true"`
	Id          any         // 主键
	LinkId      any         // 短链ID
	Code        any         // 短码（冗余，提高查询效率）
	Ip          any         // 访客 IP
	Region      any         // 访问地区（省市）
	Country     any         // 访问国家
	Referer     any         // 来源完整 URL
	RefererHost any         // 来源域名
	UserAgent   any         // User-Agent
	Device      any         // 设备类型：desktop / mobile / tablet
	Browser     any         // 浏览器
	Os          any         // 操作系统
	CreatedAt   *gtime.Time // 点击时间
}
