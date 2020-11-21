package dto

import (
	"gatewayDemo/public"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required,vail_username"`
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

func (params *AdminLoginInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"" validate:"required"`
}

type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}
