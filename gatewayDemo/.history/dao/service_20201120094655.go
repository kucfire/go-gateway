package dao

import (
	"gatewayDemo/dto"
	"gatewayDemo/public"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceInfo struct {
	Id          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"服务类型"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除; 0:否; 1:是"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (t *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}

// PageList : get the page mgs
func (t *ServiceInfo) PageList(c *gin.Context,
	tx *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := 0
	//
	offset := (param.PageNo - 1) * param.PageSize
	list := &ServiceInfo{}
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Where("is_delete = 0")
	if param.Info != "" {
		query = query.Where(
			"(service_name like %?% or service_desc like %?%)",
			param.Info, param.Info)
	}
	if err := query.Limit(param.PageSize).Offset(offset).Find(&list).Error; err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

}