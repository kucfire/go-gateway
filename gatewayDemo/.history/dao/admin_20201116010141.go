package dao

import (
	"gatewayDemo/dto"
	"gatewayDemo/public"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type AdminInfo struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"username" gorm:"column:username" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:""`
}

func (ad *AdminInfo) TableName() string {
	return "gateway_admin"
}

func (ad *AdminInfo) Find(c *gin.Context, tx *gorm.DB, search *AdminInfo) (*AdminInfo, error) {
	out := &AdminInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}

	return out, nil

}

func (ad *AdminInfo) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*AdminInfo, error) {
	return nil, nil
}
