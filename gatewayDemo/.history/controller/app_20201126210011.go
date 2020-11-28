package controller

import (
	"gatewayDemo/dao"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type AppController struct{}

func AppRegister(group *gin.RouterGroup) {
	app := AppController{}

	// app group
	group.GET("/app_list", app.AppList)
	group.POST("/app_add", app.AppAdd)
	// group.GET("/app_delete", app.AppDelete)
	// group.GET("/app_detail", app.AppDetail)
	// group.GET("/app_stat", app.AppStat)
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
// @ID /app/app_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_no query int true "页数"
// @Param page_size query int true "每页条数"
// @Success 200 {object} middleware.Response{data=dto.AppListInput} "success"
// @Router /app/app_list [get]
func (app *AppController) AppList(c *gin.Context) {
	params := &dto.AppListInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 从DB中分页读取基本信息
	appInfo := &dao.AppInfo{}
	list, total, err := appInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 格式化基本信息
	outList := []dto.AppListItemOutput{}
	for _, listItem := range list {

		outItem := dto.AppListItemOutput{
			ID:       listItem.ID,
			AppID:    listItem.AppID,
			Name:     listItem.Name,
			Sercet:   listItem.Secret,
			WhiteIPS: listItem.WhiteIps,
			QPD:      0,
			QPS:      0,
		}
		outList = append(outList, outItem)
	}

	out := &dto.AppListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(c, out)
}

// AppAdd godoc
// @Summary 添加租户
// @Description 添加租户
// @Tags 租户管理
// @ID /app/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto. true "body"
// @Success 200 {object} middleware.Response{data=dto.} "success"
// @Router /app/service_add_http [post]
func (app *AppController) AppAdd(c *gin.Context) {
	params := &dto.AppListInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
}

/*
-----------------------------------------------------------------------
----------------------------block tail---------------------------------
-----------------------------------------------------------------------
*/
