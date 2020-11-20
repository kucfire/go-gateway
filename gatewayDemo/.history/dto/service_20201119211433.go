package dto

import (
	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
)

type ServiceListInput struct {
}

func (params *ServiceListInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
