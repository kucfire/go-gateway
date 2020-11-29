package controller

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"
	"gatewayDemo/public"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type DashBoardController struct{}

func DashBoardRegister(group *gin.RouterGroup) {
	dashboard := &DashBoardController{}

	group.GET("panel_group_data", dashboard.PanelGroupData)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Success 200 {object} middleware.Response{data=dto.} "success"
// @Router /dashboard/panel_group_data [get]
func (dashboard *DashBoardController) PanelGroupData(c *gin.Context) {
	params := &dto.ServiceListInput{}
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

	// 从DB中分页读取基本信息
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 格式化基本信息
	outList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		// // 1. http后缀接入 ： clusterIP + clusterPort + path
		// // 2. http域名接入 ： domain
		// // 3. tcp、grpc接入： clisterIP + serverPort
		serviceAddr := "unknow"

		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}

		// http
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHTTPS == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s",
				clusterIP,
				clusterPort,
				serviceDetail.HTTPRule.Rule)
		}

		// https
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHTTPS == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s",
				clusterIP,
				clusterSSLPort,
				serviceDetail.HTTPRule.Rule)
		}

		// http domain
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}

		// tcp
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d",
				clusterIP,
				serviceDetail.TCPRule.Port)
		}

		// grpc
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d",
				clusterIP,
				serviceDetail.GRPCRule.Port)
		}

		outItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			LoadType:    listItem.LoadType,
			ServiceAddr: serviceAddr,
			Qps:         0,
			Qpd:         0,
			TotalNode:   len(serviceDetail.LoadBalance.GetIPList()),
		}
		outList = append(outList, outItem)
	}

	out := &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(c, out)
}
