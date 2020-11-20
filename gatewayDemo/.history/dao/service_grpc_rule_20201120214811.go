package dao

import (
	"gatewayDemo/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceGRPCRule struct {
	Id             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	Port           int8   `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	HeaderTransfor string `json:"black_list" gorm:"column:black_list" description:"黑名单ip"`
}

func (t *ServiceGRPCRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (t *ServiceGRPCRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (t *ServiceGRPCRule) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
