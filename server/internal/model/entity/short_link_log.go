// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ShortLinkLog is the golang structure for table short_link_log.
type ShortLinkLog struct {
	Id          int64       `json:"id"          orm:"id"           description:"主键"`
	LinkId      int64       `json:"linkId"      orm:"link_id"      description:"短链ID"`
	Code        string      `json:"code"        orm:"code"         description:"短码（冗余，提高查询效率）"`
	Ip          string      `json:"ip"          orm:"ip"           description:"访客 IP"`
	Region      string      `json:"region"      orm:"region"       description:"访问地区（省市）"`
	Country     string      `json:"country"     orm:"country"      description:"访问国家"`
	Referer     string      `json:"referer"     orm:"referer"      description:"来源完整 URL"`
	RefererHost string      `json:"refererHost" orm:"referer_host" description:"来源域名"`
	UserAgent   string      `json:"userAgent"   orm:"user_agent"   description:"User-Agent"`
	Device      string      `json:"device"      orm:"device"       description:"设备类型：desktop / mobile / tablet"`
	Browser     string      `json:"browser"     orm:"browser"      description:"浏览器"`
	Os          string      `json:"os"          orm:"os"           description:"操作系统"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"点击时间"`
}
