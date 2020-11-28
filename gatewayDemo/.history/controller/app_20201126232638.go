package controller

import (
	"errors"
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
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body dto.AppAddInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AppAddInput} "success"
// @Router /app/app_add [post]
func (app *AppController) AppAdd(c *gin.Context) {
	params := &dto.AppAddInput{}
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

	//事务开始
	tx = tx.Begin()

	// 校验服务信息
	appInfo := &dao.AppInfo{
		AppID: params.AppID,
		Name:  params.Name,
	}
	if _, err = appInfo.Find(c, tx, appInfo); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("服务已存在"))
		return
	}

	// 存储服务信息
	appModel := &dao.AppInfo{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIps: params.WhiteIPS,
	}
	// 写入数据的时候serviceModel也会更新
	if err = appModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	// 存储HTTPRule
	HTTPRule := &dao.ServiceHTTPRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHTTPS:      params.NeedHTTPS,
		NeedStripURI:   params.NeedStripURI,
		NeedWEBSocket:  params.NeedWEBSocket,
		URLRewrite:     params.URLRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err = HTTPRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	// 存储控制信息
	ServiceAccessControl := &dao.ServiceAccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err = ServiceAccessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	// 存储负载均衡信息
	serviceLoadBalance := &dao.ServiceLoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IPList:                 params.IPList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err = serviceLoadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	// 提交事务
	tx.Commit()

	middleware.ResponseSuccess(c, "app msg add successful")
}

/*
-----------------------------------------------------------------------
----------------------------block tail---------------------------------
-----------------------------------------------------------------------
*/
