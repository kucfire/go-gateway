package dao

import (
	"gatewayDemo/public"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceAccessControl struct {
	Id         int64     `json:"id" gorm:"primary_key"`
	ServiceID  int64     `json:"service_id" gorm:"column:service_id" description:"服务id"`
	openAuth   int       `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	black_list string    `json:"black_list" gorm:"column:black_list" description:"黑名单"`
	white_list string    `json:"white_list" gorm:"column:white_list" description:"白名单"`
	UpdatedAt  time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt  time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete   int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除; 0:否; 1:是"`
}

func (t *ServiceAccessControl) TableName() string {
	return "gateway_service_access_control"
}

func (t *ServiceAccessControl) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (t *ServiceAccessControl) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
