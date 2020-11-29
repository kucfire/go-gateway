package controller

import (
	"errors"
	"gatewayDemo/dao"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"
	"gatewayDemo/public"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type DashBoardController struct{}

func DashBoardRegister(group *gin.RouterGroup) {
	dashboard := &DashBoardController{}

	group.GET("panel_group_data", dashboard.PanelGroupData)
	group.GET("flow_stat", dashboard.FlowStat)
	group.GET("service_stat", dashboard.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (dashboard *DashBoardController) PanelGroupData(c *gin.Context) {
	// 连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{}
	_, totalService, err := serviceInfo.PageList(c, tx, &dto.ServiceListInput{PageNo: 0, PageSize: 10})
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	appInfo := &dao.AppInfo{}
	_, totalApp, err := appInfo.PageList(c, tx, &dto.AppListInput{PageNo: 0, PageSize: 10})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      totalService,
		AppNum:          totalApp,
		CurrentQPS:      0,
		TodayRequestNum: 0,
	}
	middleware.ResponseSuccess(c, out)
}

// FlowStat godoc
// @Summary 流量统计
// @Description 流量统计
// @Tags 大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (dashboard *DashBoardController) FlowStat(c *gin.Context) {
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

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (dashboard *DashBoardController) ServiceStat(c *gin.Context) {
	// 连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 获取服务信息
	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2002, errors.New("LoadType not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, item.Name)
	}

	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}
