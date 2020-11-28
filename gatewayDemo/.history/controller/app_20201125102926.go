package controller

import "github.com/gin-gonic/gin"

type AppController struct{}

func AppRegister(group *gin.RouterGroup) {
	app := AppController{}

	// app group
	group.GET("/app_list", app.AppList)
	// group.GET("/app_delete", app.AppDelete)
	// group.GET("/app_detail", app.AppDetail)
	// group.GET("/app_stat", app.AppStat)
	// group.POST("/app_add", app.AppAdd)
	// group.POST("/app_update", app.AppUpdat)
}

/*
-----------------------------------------------------------------------
----------------------------code block---------------------------------
----------------------------app module---------------------------------
-----------------------------------------------------------------------
*/

// AppList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /service/app_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_no query int true "页数"
// @Param page_size query int true "每页条数"
// @Success 200 {object} middleware.Response{data=dto.AppListInput} "success"
// @Router /service/service_list [get]
func (app *AppController) AppList(c *gin.Context) {

}

/*
-----------------------------------------------------------------------
----------------------------block tail---------------------------------
-----------------------------------------------------------------------
*/
