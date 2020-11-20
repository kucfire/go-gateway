package controller

import "github.com/gin-gonic/gin"

type ServiceController struct {
}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/info", admin.AdminInfo)
	group.POST("/changepwd", admin.ChangePWD)
}
