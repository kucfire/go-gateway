package dao

import (
	"gatewayDemo/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceHTTPRule struct {
	Id             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType       int8   `json:"rule_type" gorm:"column:rule_type" description:"是否开启权限 1=开启"`
	Rule           string `json:"rule" gorm:"column:black_list" description:"黑名单ip"`
	NeedHttps      string `json:"white_list" gorm:"column:white_list" description:"白名单ip"`
	NeedStripURL   string `json:"white_host_home" gorm:"column:white_host_home" description:"白名单主机"`
	NeedWEBSocket  int    `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit" description:"客户端ip限流"`
	URLRewrite     int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务器限流"`
	HeaderTransfor int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务器限流"`
}

func (t *ServiceHTTPRule) TableName() string {
	return "gateway_service_http_rule"
}

func (t *ServiceHTTPRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
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
