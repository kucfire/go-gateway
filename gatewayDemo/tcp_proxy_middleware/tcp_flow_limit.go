package tcp_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/public"
	"strings"

	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func TCPFLowLimitModeMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		sInterface := c.Get("service")
		if sInterface == nil {
			c.conn.Write([]byte(fmt.Sprintf("code %d err : get server empty", 4000)))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)

		// 获取到ip信息
		list := strings.Split(c.conn.RemoteAddr().String(), ":")
		if len(list) != 2 {
			c.conn.Write([]byte(fmt.Sprintf("code %d err : addr format error", 4001)))
			c.Abort()
			return
		}
		clientIp := list[0]

		// 零值校验
		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			// 服务端限流
			serviceLimit, err := public.FlowLimiterHandler.GetFlowLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				c.conn.Write([]byte(fmt.Sprintf("code %d err : %v", 5001, err)))
				c.Abort()
				return
			}
			if !serviceLimit.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("code %d err : %v", 5002, errors.New(fmt.Sprintf("service flow limit %v ", serviceDetail.AccessControl.ServiceFlowLimit)))))
				c.Abort()
				return
			}

		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			// 客户端限流
			clientLimit, err := public.FlowLimiterHandler.GetFlowLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIp,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				c.conn.Write([]byte(fmt.Sprintf("code %d err : %v", 5003, err)))
				c.Abort()
				return
			}
			if !clientLimit.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("code %d err : %v", 5004, errors.New(fmt.Sprintf("service flow limit %v ", serviceDetail.AccessControl.ClientIPFlowLimit)))))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
