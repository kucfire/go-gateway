package controller

import "github.com/gin-gonic/gin"

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	admin := &ServiceController{}
}
