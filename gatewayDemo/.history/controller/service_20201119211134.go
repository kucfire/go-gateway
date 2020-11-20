package controller

import "github.com/gin-gonic/gin"

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list")
}

func (service *ServiceController)()  {
	
}
