package dto

import (
	"gatewayDemo/public"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminInfoOutput struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avator       string    `json:"avator"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

type ChangePWDInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

func (params *AdminLoginInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
