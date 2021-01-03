package http_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/middleware"
	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func HTTPFLowLimitModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)

		// 零值校验
		if serviceDetail.AccessControl.ServiceFlowLimit > 0 {
			// 服务端限流
			serviceLimit, err := public.FlowLimiterHandler.GetFlowLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !serviceLimit.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit)))
				c.Abort()
				return
			}

		}

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			// 客户端限流
			clientLimit, err := public.FlowLimiterHandler.GetFlowLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+c.ClientIP(),
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !clientLimit.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf("client ip flow limit %v", serviceDetail.AccessControl.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
