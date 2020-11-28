package controller

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gatewayDemo/dao"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type AppController struct{}

func AppRegister(group *gin.RouterGroup) {
	app := AppController{}

	// app group
	group.GET("/app_list", app.AppList)
	group.POST("/app_add", app.AppAdd)
	group.POST("/app_update", app.AppUpdate)
	group.GET("/app_delete", app.AppDelete)
	group.GET("/app_detail", app.AppDetail)
	group.GET("/app_stat", app.AppStat)
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
			QPD:      listItem.QPD,
			QPS:      listItem.QPS,
			RealQPD:  0,
			RealQPS:  0,
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
	// tx = tx.Begin()

	// 校验服务信息
	appInfo := &dao.AppInfo{
		AppID: params.AppID,
	}
	if _, err = appInfo.Find(c, tx, appInfo); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("租客ID已被占用"))
		return
	}

	if params.Secret == "" {
		h := md5.New()
		h.Write([]byte(params.AppID))
		params.Secret = hex.EncodeToString(h.Sum(nil))
	}

	// 存储服务信息
	appModel := &dao.AppInfo{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIps: params.WhiteIPS,
		QPS:      params.QPS,
		QPD:      params.QPD,
	}
	// 写入数据的时候serviceModel内的参数也会更新
	if err = appModel.Save(c, tx); err != nil {
		// tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	// // 提交事务
	// tx.Commit()

	middleware.ResponseSuccess(c, "app msg add successful")
}

// AppUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags 租户管理
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto.AppUpdateInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AppUpdateInput} "success"
// @Router /app/app_update [post]
func (adminligin *AppController) AppUpdate(c *gin.Context) {
	params := &dto.AppUpdateInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
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
		ID: params.ID,
	}
	appInfo, err = appInfo.Find(c, tx, appInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	if params.Secret == "" {
		h := md5.New()
		h.Write([]byte(params.AppID))
		params.Secret = hex.EncodeToString(h.Sum(nil))
	}

	// 更新租客信息
	appUpdate := &dao.AppInfo{
		ID:        appInfo.ID,
		AppID:     appInfo.AppID,
		Name:      params.Name,
		Secret:    params.Secret,
		WhiteIps:  params.WhiteIPS,
		CreatedAt: appInfo.CreatedAt,
		IsDelete:  appInfo.IsDelete,
	}
	if err := appUpdate.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	// 提交事务
	tx.Commit()

	middleware.ResponseSuccess(c, "app msg updated successful")

}

// AppDelete godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [get]
func (service *AppController) AppDelete(c *gin.Context) {
	params := &dto.AppDeleteInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 从DB中读取基本信息
	appInfo := &dao.AppInfo{ID: params.ID}
	appInfo, err = appInfo.Find(c, tx, appInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	appInfo.IsDelete = 1
	if err = appInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "app deleted successful")
}

// AppDetail godoc
// @Summary 租户信息
// @Description 租户信息
// @Tags 租户管理
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dao.AppDeleteInput} "success"
// @Router /app/app_detail [get]
func (service *AppController) AppDetail(c *gin.Context) {
	params := &dto.AppDeleteInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 检验服务信息是否存在
	appInfo := &dao.AppInfo{ID: params.ID}
	appInfo, err = appInfo.Find(c, tx, appInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// appDetail, err := appInfo.Find(c, tx, appInfo)
	// if err != nil {
	// 	middleware.ResponseError(c, 2003, err)
	// 	return
	// }

	middleware.ResponseSuccess(c, appInfo)
}

// AppStat godoc
// @Summary 租户统计
// @Description 租户统计
// @Tags 租户管理
// @ID /app/app_stat
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dao.AppStatOutput} "success"
// @Router /app/app_stat [get]
func (service *AppController) AppStat(c *gin.Context) {
	// 由于只需要一个ID所以直接调用delete的输入结构即可
	// params := &dto.ServiceDeleteInput{}
	// if err := params.BindingValidParams(c); err != nil {
	// 	// log.F  atal("params.BindingValidParams err : %v", err)
	// 	middleware.ResponseError(c, 2000, err)
	// 	return
	// }

	// // 连接池
	// tx, err := lib.GetGormPool("default")
	// if err != nil {
	// 	middleware.ResponseError(c, 2001, err)
	// 	return
	// }

	// // 从DB中读取基本信息
	// serviceInfo := &dao.ServiceInfo{ID: params.ID}
	// serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	// if err != nil {
	// 	middleware.ResponseError(c, 2002, err)
	// 	return
	// }

	TodayList := []int64{}
	for i := 0; i < time.Now().Hour(); i++ {
		TodayList = append(TodayList, 0)
	}

	YesterdayList := []int64{}
	for i := 0; i < 24; i++ {
		YesterdayList = append(YesterdayList, 0)
	}

	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Yesterday: YesterdayList,
		Today:     TodayList,
	})
}

/*
-----------------------------------------------------------------------
----------------------------block tail---------------------------------
-----------------------------------------------------------------------
*/
