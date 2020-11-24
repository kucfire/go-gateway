package dao

import (
	"gatewayDemo/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceHTTPRule struct {
	Id             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 0=url前缀url_prefix 1=域名domain"`
	Rule           string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix时标是url前缀"`
	NeedHttps      int8   `json:"need_https" gorm:"column:need_https" description:"支持https 1=支持"`
	NeedStripURL   int8   `json:"need_strip_url" gorm:"column:need_strip_url" description:"启用strip_url 1=启用"`
	NeedWEBSocket  int8   `json:"need_websocket" gorm:"column:need_websocket" description:"是否支持websocket 1=支持"`
	URLRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔"`
}

func (t *ServiceHTTPRule) TableName() string {
	return "gateway_service_http_rule"
}

func (t *ServiceHTTPRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceHTTPRule) (*ServiceHTTPRule, error) {
	out := &ServiceHTTPRule{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (t *ServiceHTTPRule) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
