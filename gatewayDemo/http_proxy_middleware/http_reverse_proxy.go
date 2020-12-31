package http_proxy_middleware

import (
	"gatewayDemo/dao"
	"gatewayDemo/middleware"
	"gatewayDemo/reverse_proxy/proxy"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否已经将对应的serviceDetail存放进gin的上下文
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail) // 接口转换

		// 获取对应的负载均衡设置
		lb, err := dao.LoadBalanceHandler.GetLoadBalance(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 1003, err)
			c.Abort()
			return
		}

		// 获取transport
		transport, err := dao.TransportorHandler.GetTransportor(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 1004, err)
			c.Abort()
			return
		}

		// 创建 reverseproxy
		// 使用 reverseproxy.ServerHTTP(c.request, c.Response)
		proxy := proxy.NewLoadBalanceReverseProxy(c, lb, transport)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}
}
