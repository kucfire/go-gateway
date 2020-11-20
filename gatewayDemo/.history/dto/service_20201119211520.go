package dto

import (
	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
)

type ServiceListInput struct {
	// 关键词
	Info string `json:"info" form:"info" comment:"关键词" example:"" validate:""`
}

func (params *ServiceListInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
