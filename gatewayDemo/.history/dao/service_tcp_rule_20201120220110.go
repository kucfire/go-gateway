package dao

import (
	"gatewayDemo/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceTCPRule struct {
}

func (t *ServiceTCPRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (t *ServiceTCPRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (t *ServiceTCPRule) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
