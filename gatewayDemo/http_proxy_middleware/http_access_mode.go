package http_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/middleware"

	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
)

// 匹配接入方式 给予请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			// 阻止后续的处理函数
			c.Abort()
			return
		}

		fmt.Printf("matched service : %v\n", public.ObjToJson(serviceDetail))
		c.Set("serviceDetail", serviceDetail)
		c.Next()
	}
}
