package controller

import (
	"errors"
	"gatewayDemo/dao"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"
	"strings"

	"github.com/e421083458/golang_common/lib"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type AppController struct{}

func AppRegister(group *gin.RouterGroup) {
	app := AppController{}

	// app group
	group.GET("/app_list", app.AppList)
	group.POST("/app_add", app.AppAdd)
	group.POST("/app_update", app.AppUpdate)
	// group.GET("/app_delete", app.AppDelete)
	// group.GET("/app_detail", app.AppDetail)
	// group.GET("/app_stat", app.AppStat)
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
	// 写入数据的时候serviceModel内的参数也会更新
	if err = appModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	// 提交事务
	tx.Commit()

	middleware.ResponseSuccess(c, "app msg add successful")
}

// AppUpdate godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto. true "body"
// @Success 200 {object} middleware.Response{data=dto.} "success"
// @Router /app/app_update [post]
func (adminligin *AppController) AppUpdate(c *gin.Context) {
	params := &dto.{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	//校验ip列表和权重列表
	if len(strings.Split(params.IPList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		middleware.ResponseError(c, 2001, errors.New("IPList与权重列表数量不一致"))
		return
	}

	// 连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//事务开始
	tx = tx.Begin()

	// 检验服务信息是否存在
	serviceInfoSearch := &dao.ServiceInfo{ID: params.ID}
	serviceInfoSearch, err = serviceInfoSearch.Find(c, tx, serviceInfoSearch)
	if err != nil && err == gorm.ErrRecordNotFound {
		middleware.ResponseError(c, 2003, errors.New("服务不存在"))
		return
	}

	// 校验服务信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	// 更新服务信息
	serviceInfo = &dao.ServiceInfo{
		ID:          serviceInfoSearch.ID,
		LoadType:    serviceInfoSearch.LoadType,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
		CreatedAt:   serviceInfoSearch.CreatedAt,
		IsDelete:    serviceInfoSearch.IsDelete,
	}
	if err := serviceInfo.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	// 存储HTTPRule
	serviceHTTPRule := serviceDetail.HTTPRule
	serviceHTTPRule.RuleType = params.RuleType
	serviceHTTPRule.Rule = params.Rule
	serviceHTTPRule.NeedHTTPS = params.NeedHTTPS
	serviceHTTPRule.NeedStripURI = params.NeedStripURI
	serviceHTTPRule.NeedWEBSocket = params.NeedWEBSocket
	serviceHTTPRule.URLRewrite = params.URLRewrite
	serviceHTTPRule.HeaderTransfor = params.HeaderTransfor
	if err = serviceHTTPRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	// 存储控制信息
	serviceAccessControl := serviceDetail.AccessControl
	serviceAccessControl.OpenAuth = params.OpenAuth
	serviceAccessControl.BlackList = params.BlackList
	serviceAccessControl.WhiteList = params.WhiteList
	serviceAccessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	serviceAccessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err = serviceAccessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	// 存储负载均衡信息
	serviceLoadBalance := serviceDetail.LoadBalance
	serviceLoadBalance.RoundType = params.RoundType
	serviceLoadBalance.IPList = params.IPList
	serviceLoadBalance.WeightList = params.WeightList
	serviceLoadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	serviceLoadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	serviceLoadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	serviceLoadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err = serviceLoadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	// 提交事务
	tx.Commit()

	middleware.ResponseSuccess(c, "HTTP msg updated successful")

}

/*
-----------------------------------------------------------------------
----------------------------block tail---------------------------------
-----------------------------------------------------------------------
*/
