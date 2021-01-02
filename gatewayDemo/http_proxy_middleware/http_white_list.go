package http_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/middleware"
	"gatewayDemo/public"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func HTTPWhiteListModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)

		ipList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, c.ClientIP()) {
				middleware.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
