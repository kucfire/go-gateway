package dao

import "time"

type AppInfo struct {
	ID        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"租户id"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIps  string    `json:"white_ips" gorm:"column:white_ips" description:"密钥"`
	QPD       string    `json:"qpd" gorm:"column:qpd" description:"密钥"`
	QPS       string    `json:"qps" gorm:"column:qps" description:"密钥"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}
