package http_proxy_middleware

import (
	"fmt"
	"gatewayDemo/dao"
	"gatewayDemo/middleware"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 匹配接入方式 给予请求信息
func HTTPURLRewriteModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)

		for _, item := range strings.Split(serviceDetail.HTTPRule.URLRewrite, ",") {
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}

			regexp, err := regexp.Compile(items[0])
			if err != nil {
				fmt.Println("regexp.Compile err : ", err)
				continue
			}
			replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(replacePath)
		}

		c.Next()
	}
}
