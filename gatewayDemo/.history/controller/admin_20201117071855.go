package controller

import "github.com/gin-gonic/gin"

type AdminController struct {
}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/admininfo", admin.)
}
