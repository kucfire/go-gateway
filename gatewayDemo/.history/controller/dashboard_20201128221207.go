package controller

import "github.com/gin-gonic/gin"

type DashBoardController struct{}

func DashBoardRegister(group *gin.RouterGroup) {
	service := &ServiceController{}

	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_stat", service.ServiceStat)
}
