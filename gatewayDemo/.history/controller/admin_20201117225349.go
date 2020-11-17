package controller

import (
	"gatewayDemo/dto"
	"gatewayDemo/middleware"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/admin_info", admin.AdminInfo)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginInput} "success"
// @Router /admin/admin_info [get]
func (admininfo *AdminController) AdminInfo(c *gin.Context) {
	// 1. 读取sessionKey对应json，转换为结构体
	// 2. 取出数据然后封装输出结构体

	out := &dto.AdminLoginOutput{}
	middleware.ResponseSuccess(c, out)
}
