package dao

import "time"

type ServiceInfo struct {
	Id          int       `json:"id" gorm:"primary_key" description:"自增主键"`
	LoadType    string    `json:"username" gorm:"column:user_name" description:"管理员用户名"`
	ServiceName string    `json:"salt" gorm:"column:salt" description:"盐"`
	ServiceDesc string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}
