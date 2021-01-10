package tcp_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/public"
	"strings"
)

// 匹配接入方式 给予请求信息
func TCPBlackListModeMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		sInterface := c.Get("service")
		if sInterface == nil {
			c.conn.Write([]byte(fmt.Sprintf("code %d err : get server empty", 4000)))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)

		blackIpList := []string{}
		whiteIpList := []string{}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}
		if serviceDetail.AccessControl.WhiteList != "" {
			whiteIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		// 获取到ip信息
		list := strings.Split(c.conn.RemoteAddr().String(), ":")
		if len(list) != 2 {
			c.conn.Write([]byte(fmt.Sprintf("code %d err : addr format error", 4001)))
			c.Abort()
			return
		}
		clientIp := list[0]
		if serviceDetail.AccessControl.OpenAuth == 1 && len(blackIpList) > 0 && len(whiteIpList) == 0 {
			if public.InStringSlice(blackIpList, clientIp) {
				c.conn.Write([]byte(fmt.Sprintf("code %d err : %s in black ip list", 4001, clientIp)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
