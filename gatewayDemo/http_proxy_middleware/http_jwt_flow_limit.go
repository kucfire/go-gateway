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
func HTTPFJwtLowLimitModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("appInfo")
		if !ok {
			middleware.ResponseError(c, 5000, errors.New("AppInfo is not find"))
			c.Abort()
			return
		}
		appInfo := appInterface.(*dao.AppInfo)

		if appInfo.QPS > 0 {
			clientLimit, err := public.FlowLimiterHandler.GetFlowLimiter(
				public.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.QPS))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}

			if !clientLimit.Allow() {
				middleware.ResponseError(c, 5002,
					errors.New(fmt.Sprintf("app client ip:%v flow limit %v", c.ClientIP(), appInfo.QPS)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
