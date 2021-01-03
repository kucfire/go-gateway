package http_proxy_middleware

import (
	"gatewayDemo/dao"
	"gatewayDemo/middleware"
	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func HTTPFLowCountModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)

		// 统计项
		// 1. 全站统计
		// 2. 服务统计
		// 3. 租户统计
		// 全站统计
		totalCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowTotalPrefix)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		totalCounter.Increase()
		// dayTotalCount, _ := totalCounter.GetDayData(time.Now())
		// fmt.Printf("{%s totalCounter —— qps:%v total:%v}\n", serviceDetail.Info.ServiceName, totalCounter.QPS, dayTotalCount)
		// 服务统计
		serviceCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		serviceCounter.Increase()
		// dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		// fmt.Printf("{%s serviceCount —— qps:%v total:%v}\n", serviceDetail.Info.ServiceName, serviceCounter.QPS, dayServiceCount)

		c.Next()
	}
}
