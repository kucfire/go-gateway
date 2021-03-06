package http_proxy_middleware

import (
	"gatewayDemo/dao"
	"gatewayDemo/middleware"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func HTTPHeaderTransferModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)
		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}
			if items[0] == "add" || items[0] == "edit" {
				c.Request.Header.Set(items[1], items[2])
			}

			if items[0] == "del" {
				c.Request.Header.Del(items[1])
			}
		}

		c.Next()
	}
}
