package controller

import "github.com/gin-gonic/gin"

type DashBoardController struct{}

func DashBoardRegister(group *gin.RouterGroup) {
	dashboard := &DashBoardController{}

	group.GET("panel_group_data", dashboard)
}
