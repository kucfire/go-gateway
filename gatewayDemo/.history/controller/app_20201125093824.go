package controller

import "github.com/gin-gonic/gin"

type AppController struct{}

func AppRegister(group *gin.RouterGroup) {
	app := AppController{}

	// app message
	group.GET("/service_list", app.AppList)
	group.GET("/service_delete", app.AppDelete)
	group.GET("/service_detail", app.AppDetail)
	group.GET("/service_stat", app.AppStat)

	// app group
	group.POST("/service_add_http", app.AppAdd)
	group.POST("/service_update_http", app.AppUpdatApp)
}
