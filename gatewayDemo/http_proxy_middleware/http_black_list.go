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
func HTTPBlackListModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
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
		if serviceDetail.AccessControl.OpenAuth == 1 && len(blackIpList) > 0 && len(whiteIpList) == 0 {
			if public.InStringSlice(blackIpList, c.ClientIP()) {
				middleware.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s in black ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
