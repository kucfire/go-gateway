package controller

import (
	"encoding/base64"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type OauthController struct{}

func OauthRegister(group *gin.RouterGroup) {
	Oauth := &OauthController{}
	// login interface
	group.POST("/tokens", Oauth.Token)

}

// Token godoc
// @Summary 获取Token
// @Description 获取Token
// @Tags Token接口
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (oauth *OauthController) Token(c *gin.Context) {
	params := &dto.TokensInput{}
	// 绑定参数进上下文
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 获取header里面的Authorization
	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		middleware.ResponseError(c, 2001, errors.New("Authorization's format is error/用户名或密码格式错误"))
		return
	}

	appScore, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// fmt.Println("Authorization : ", string(appScore)) // 测试数据

	// 取出 app_id secret
	// 生成app_list
	// 匹配app_id
	// 基于jwt生成token
	// 生成output
	parts := strings.Split(string(appScore), ":")
	if len(parts) != 2 {
		middleware.ResponseError(c, 2003, errors.New("Authorization's format is error/用户名或密码格式错误"))
		return
	}
	// fmt.Println("Authorization : ", parts) // 测试数据
	appID := parts[0]
	secret := parts[1]

	out := &dto.TokensOutput{}
	middleware.ResponseSuccess(c, out)

}
