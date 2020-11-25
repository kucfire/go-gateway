package controller

import "github.com/gin-gonic/gin"

type AppController struct{}

func AppRegister(group *gin.RouterGroup) {
	app := AppController{}

	group.GET("/service_list", app.ServiceList)
	group.GET("/service_delete", app.ServiceDelete)
	group.GET("/service_detail", app.ServiceDetail)
	group.GET("/service_stat", app.ServiceStat)

	// app group
	group.POST("/service_add_http", app.ServiceAddHTTP)
	group.POST("/service_update_http", app.ServiceUpdateHTTP)
}
