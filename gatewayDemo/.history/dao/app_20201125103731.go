package dao

import (
	"gatewayDemo/dto"
	"gatewayDemo/public"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// AppInfo : 租户信息结构体
type AppInfo struct {
	ID        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"租户id"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIps  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	QPD       string    `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	QPS       string    `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

// TableName : 获取数据库对应的表名称
func (ad *AppInfo) TableName() string {
	return "gateway_app"
}

// Find : 查询
func (ad *AppInfo) Find(c *gin.Context, tx *gorm.DB, search *AppInfo) (*AppInfo, error) {
	out := &AppInfo{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Save : 存储信息
func (ad *AppInfo) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(ad).Error
}

// PageList : 获取相应页数的信息
func (ad *AppInfo) PageList(c *gin.Context,
	tx *gorm.DB, param *dto.AppListInput) ([]AppInfo, int64, error) {
	total := int64(0)
	// 从第几条开始查询
	offset := (param.PageNo - 1) * param.PageSize
	list := []AppInfo{}

	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(ad.TableName()).Where("is_delete = 0")
	if param.Info != "" {
		query = query.Where(
			"(name like ?)",
			"%"+param.Info+"%")
	}
	// Limit : 限制查询多少条结果; offset : 从第几条开始查询
	if err := query.Limit(param.PageSize).Offset(offset).Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)

	return list, total, nil
}
