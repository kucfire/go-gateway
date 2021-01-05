package dto

import (
	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
)

type TokensInput struct {
	GrantType string `json:"grant_type" form:"grant_type" comment:"授权类型" example:"client_credentials" validate:"required"`
	Scope     string `json:"scope" form:"scope" comment:"权限范围" example:"read_write" validate:"required"`
}

// BindingValidParams ： 参数验证
func (params *TokensInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type TokensOutput struct {
	AccessToken string `json:"access_token" form:"access_token"`
	Expires     int    `json:"expires" form:"expires"`
	TokenType   string `json:"token_type" form:"token_type"`
	Scope       string `json:"scope" form:"scope"`
}
