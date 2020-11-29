package controller

import "github.com/gin-gonic/gin"

type DashBoardController struct{}

func DashBoardRegister(group *gin.RouterGroup) {
	service := &ServiceController{}

	group.GET("/service_list")
}
