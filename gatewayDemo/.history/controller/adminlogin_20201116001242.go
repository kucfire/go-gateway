package controller

import (
	"gatewayDemo/dto"
	"gatewayDemo/middleware"

	"github.com/gin-gonic/gin"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)

}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginInput} "success"
// @Router /admin_login/login [post]
func (adminligin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 1001, err)
		return
	}

	// 1. 从数据库中取得管理员信息 adminInfo

	out := &dto.AdminLoginOutput{Token: params.Username}
	middleware.ResponseSuccess(c, out)

}