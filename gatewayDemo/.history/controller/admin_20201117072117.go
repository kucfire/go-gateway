package controller

import "github.com/gin-gonic/gin"

type AdminController struct {
}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/admininfo", admin.AdminInfo)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admininfo/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginInput} "success"
// @Router /admininfo/login [post]
func (admininfo *AdminController) AdminInfo(c *gin.Context) {

}
