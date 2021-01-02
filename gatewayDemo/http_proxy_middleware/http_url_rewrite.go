package http_proxy_middleware

import (
	"gatewayDemo/dao"
	"gatewayDemo/middleware"
	"gatewayDemo/public"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func HTTPStripURIModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)

		// http://127.0.0.1:8080/http_test_string/aaa
		// http://127.0.0.1:2003/aaa
		if serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripURI == 1 {
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule, "", 1)
		}

		c.Next()
	}
}
