package dao

import (
	"gatewayDemo/public"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

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

func (ad *AppInfo) Find(c *gin.Context, tx *gorm.DB, search *AppInfo) (*AppInfo, error) {
	out := &AppInfo{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (ad *AppInfo) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(ad).Error
}
