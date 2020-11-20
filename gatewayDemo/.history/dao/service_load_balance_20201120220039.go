package dao

import (
	"gatewayDemo/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceLoadBalance struct {
	Id                     int64  `json:"id" gorm:"primary_key"`
	ServiceID              int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	CheckMethod            int    `json:"check_method" gorm:"column:check_method" description:"匹配类型 0=url前缀url_prefix 1=域名domain"`
	CheckTimeout           int    `json:"check_timeout" gorm:"column:check_timeout" description:"type=domain表示域名，type=url_prefix时标是url前缀"`
	CheckInterval          int    `json:"check_interval" gorm:"column:check_interval" description:"支持https 1=支持"`
	RoundType              int8   `json:"round_type" gorm:"column:round_type" description:"启用strip_url 1=启用"`
	IPList                 string `json:"ip_list" gorm:"column:ip_list" description:"是否支持websocket 1=支持"`
	WeightList             string `json:"weight_list" gorm:"column:weight_list" description:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔"`
	ForbidList             string `json:"forbid_list" gorm:"column:forbid_list" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔"`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔"`
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔"`
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔"`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔"`
}

func (t *ServiceLoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *ServiceLoadBalance) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (t *ServiceLoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
