package dao

import (
	"fmt"
	"gatewayDemo/dto"
	"gatewayDemo/public"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type AdminInfo struct {
	ID        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"username" gorm:"column:user_name" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (ad *AdminInfo) TableName() string {
	return "gateway_admin"
}

func (ad *AdminInfo) Find(c *gin.Context, tx *gorm.DB, search *AdminInfo) (*AdminInfo, error) {
	out := &AdminInfo{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (ad *AdminInfo) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*AdminInfo, error) {
	admin, err := ad.Find(c, tx, (&AdminInfo{UserName: param.UserName, IsDelete: 0}))
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}

	saltPassword := public.GenSaltPassword(admin.Salt, param.Password)
	fmt.Println(saltPassword, admin.Password)
	if saltPassword != admin.Password {
		return nil, errors.New("密码错误，请重新输入")
	}
	return admin, nil
}

func (ad *AdminInfo) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(ad).Error
}
