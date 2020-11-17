package controller

import "github.com/gin-gonic/gin"

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
// @Router /admin/admin_info [post]
func (admininfo *AdminController) AdminInfo(c *gin.Context) {

}
