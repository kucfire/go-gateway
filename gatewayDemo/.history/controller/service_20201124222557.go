package controller

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/dto"
	"gatewayDemo/middleware"
	"gatewayDemo/public"
	"strings"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}

	// HTTP group
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_stat", service.ServiceStat)

	// GRPC group
	group.POST("/service_add_grpc", service.ServiceAddGRPC)
	group.POST("/service_update_grpc", service.ServiceUpdateGRPC)

	// TCP group
	group.POST("/service_add_grpc", service.ServiceAddTCP)
	group.POST("/service_update_grpc", service.ServiceUpdateTCP)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_no query int true "页数"
// @Param page_size query int true "每页条数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListInput} "success"
// @Router /service/service_list [get]
func (service *ServiceController) ServiceList(c *gin.Context) {
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

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
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
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceInfo.IsDelete = 1
	if err = serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "deleted successful")
}

// ServiceDetail godoc
// @Summary 服务信息
// @Description 服务信息
// @Tags 服务管理
// @ID /service/service_detail
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router /service/service_detail [get]
func (service *ServiceController) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
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
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, serviceDetail)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 服务管理
// @ID /service/service_stat
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceStatOutput} "success"
// @Router /service/service_stat [get]
func (service *ServiceController) ServiceStat(c *gin.Context) {
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
	------------------------HTTP MODULE----------------------------
*/

// ServiceAddHTTP godoc
// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务管理
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=dto.ServiceAddHTTPInput} "success"
// @Router /service/service_add_http [post]
func (adminligin *ServiceController) ServiceAddHTTP(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	//校验ip列表和权重列表
	if len(strings.Split(params.IPList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		middleware.ResponseError(c, 2004, errors.New("IPList与权重列表数量不一致"))
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
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if _, err = serviceInfo.Find(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("服务已存在"))
		return
	}

	// 校验域名
	HTTPURL := &dao.ServiceHTTPRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err = HTTPURL.Find(c, tx, HTTPURL); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("服务接入前缀或域名已存在"))
		return
	}

	// 存储服务信息
	serviceModel := &dao.ServiceInfo{
		LoadType:    public.LoadTypeHTTP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	// 写入数据的时候serviceModel也会更新
	if err = serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
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

	middleware.ResponseSuccess(c, "HTTP msg add successful")

}

// ServiceUpdateHTTP godoc
// @Summary 修改HTTP服务
// @Description 修改HTTP服务
// @Tags 服务管理
// @ID /service/service_update_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=dto.ServiceUpdateHTTPInput} "success"
// @Router /service/service_update_http [post]
func (adminligin *ServiceController) ServiceUpdateHTTP(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
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

	// 校验服务信息
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("服务已存在"))
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
		middleware.ResponseError(c, 2004, err)
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
		middleware.ResponseError(c, 2005, err)
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
		middleware.ResponseError(c, 2006, err)
		return
	}

	// 提交事务
	tx.Commit()

	middleware.ResponseSuccess(c, "HTTP msg updated successful")

}

/*
	------------------------GRPC MODULE----------------------------
*/

// ServiceAddGRPC godoc
// @Summary 添加GRPC服务
// @Description 添加GRPC服务
// @Tags 服务管理
// @ID /service/service_add_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddGRPCInput true "body"
// @Success 200 {object} middleware.Response{data=dto.ServiceAddHTTPInput} "success"
// @Router /service/service_add_grpc [post]
func (adminligin *ServiceController) ServiceAddGRPC(c *gin.Context) {
	params := &dto.ServiceAddGRPCInput{}
	if err := params.BindingValidParams(c); err != nil {
		// log.F  atal("params.BindingValidParams err : %v", err)
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 校验 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(c, lib.GORMDefaultPool, infoSearch); err == nil {
		middleware.ResponseError(c, 2001, errors.New("服务已存在"))
		return
	}

	// 校验端口是否被占用,需同时检测tcp和grpc两边的port,避免冲突
	// 检验tcp的port
	tcpRuleSearch := &dao.ServiceTCPRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(c, lib.GORMDefaultPool, tcpRuleSearch); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务端口被占用"))
		return
	}
	// 校验grpc的port
	grpcRuleSearch := &dao.ServiceGRPCRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(c, lib.GORMDefaultPool, grpcRuleSearch); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用"))
		return
	}

	//校验ip列表和权重列表
	if len(strings.Split(params.IPList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		middleware.ResponseError(c, 2004, errors.New("IPList与权重列表数量不一致"))
		return
	}

	// 连接池,事物开始
	tx := lib.GORMDefaultPool.Begin()

	// 存储服务信息
	serviceModel := &dao.ServiceInfo{
		LoadType:    public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	// 写入数据的时候serviceModel也会更新
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	// 存储负载均衡信息
	serviceLoadBalance := &dao.ServiceLoadBalance{
		ServiceID:  serviceModel.ID,
		RoundType:  params.RoundType,
		IPList:     params.IPList,
		WeightList: params.WeightList,
	}
	if err := serviceLoadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	// 存储GRPCRule
	grpcRule := &dao.ServiceGRPCRule{
		ServiceID:      infoSearch.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
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
	if err := ServiceAccessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	// 提交事务
	tx.Commit()

	middleware.ResponseSuccess(c, "GRPC msg add successful")

}

// ServiceUpdateGRPC godoc
// @Summary 修改GRPC服务
// @Description 修改GRPC服务
// @Tags 服务管理
// @ID /service/service_update_grpc
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateGRPCInput true "body"
// @Success 200 {object} middleware.Response{data=dto.ServiceUpdateGRPCInput} "success"
// @Router /service/service_update_grpc [post]
func (adminligin *ServiceController) ServiceUpdateGRPC(c *gin.Context) {
	params := &dto.ServiceUpdateGRPCInput{}
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

	// 连接池,事物开始
	tx := lib.GORMDefaultPool.Begin()

	// 校验服务信息
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("服务已存在"))
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
		middleware.ResponseError(c, 2004, err)
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
		middleware.ResponseError(c, 2005, err)
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
		middleware.ResponseError(c, 2006, err)
		return
	}

	// 提交事务
	tx.Commit()

	middleware.ResponseSuccess(c, "HTTP msg updated successful")

}
