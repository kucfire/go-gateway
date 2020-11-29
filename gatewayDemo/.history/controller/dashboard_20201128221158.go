package controller

import "github.com/gin-gonic/gin"

type DashBoardController struct{}

func DashBoardRegister(group *gin.RouterGroup) {
	service := &ServiceController{}

	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_stat", service.ServiceStat)

	// HTTP group
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)

	// GRPC group
	group.POST("/service_add_grpc", service.ServiceAddGRPC)
	group.POST("/service_update_grpc", service.ServiceUpdateGRPC)

	// TCP group
	group.POST("/service_add_tcp", service.ServiceAddTCP)
	group.POST("/service_update_tcp", service.ServiceUpdateTCP)
}
