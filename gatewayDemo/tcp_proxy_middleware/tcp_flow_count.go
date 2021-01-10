package tcp_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/public"
)

// 匹配接入方式 给予请求信息
func TCPFLowCountModeMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		sInterface := c.Get("service")
		if sInterface == nil {
			c.conn.Write([]byte(fmt.Sprintf("code %d err : get server empty", 4000)))
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
			c.conn.Write([]byte(fmt.Sprintf("code %d err : %v", 4001, err)))
			c.Abort()
			return
		}
		totalCounter.Increase()
		// 服务统计
		serviceCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			c.conn.Write([]byte(fmt.Sprintf("code %d err : %v", 4002, err)))
			c.Abort()
			return
		}
		serviceCounter.Increase()

		c.Next()
	}
}
