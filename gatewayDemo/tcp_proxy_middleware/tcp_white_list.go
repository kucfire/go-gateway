package tcp_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/public"
	"strings"

	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func TCPWhiteListModeMiddleware() func(c *TcpSliceRouterContext) {
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
			c.conn.Write([]byte(fmt.Sprintf("code %d err : addr format error", 3000)))
			c.Abort()
			return
		}
		clientIp := list[0]

		ipList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, clientIp) {
				// middleware.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s not in white ip list", clientIp)))
				c.conn.Write([]byte(fmt.Sprintf("code %d err : %v", 3001, errors.New(fmt.Sprintf("%s not in white ip list", clientIp)))))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
