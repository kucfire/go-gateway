package controller

import (
	"gatewayDemo/dto"
	"gatewayDemo/middleware"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list")
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param body body string true "关键词"
// @Param body body string true "页数"
// @Param body body string true "每页条数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListInput} "success"
// @Router /service/service_list [get]
func (service *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	out := &dto.ServiceListOutput{}
	middleware.ResponseSuccess(c, out)
}
