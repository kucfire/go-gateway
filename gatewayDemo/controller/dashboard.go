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

	totalCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowTotalPrefix)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      totalService,
		AppNum:          totalApp,
		CurrentQPS:      totalCounter.QPS,
		TodayRequestNum: totalCounter.TotalCount,
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

	counter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowTotalPrefix)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	TodayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		TodayList = append(TodayList, hourData)
	}

	YesterdayList := []int64{}
	yesterdayTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i < 24; i++ {
		dateTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		YesterdayList = append(YesterdayList, hourData)
	}

	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Yesterday: YesterdayList,
		Today:     TodayList,
	})
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 大盘rest
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
	result := []dto.DashServiceStatListOutput2{}
	for _, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2002, errors.New("LoadType not found"))
			return
		}
		// list[index].Name = name
		result = append(result, dto.DashServiceStatListOutput2{
			Name:  name,
			Value: item.Value,
		})
		legend = append(legend, name)
	}

	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   result,
	}
	middleware.ResponseSuccess(c, out)
}
