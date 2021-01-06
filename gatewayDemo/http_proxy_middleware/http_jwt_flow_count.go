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
func HTTPJwtFLowCountModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// sInterface, ok := c.Get("serviceDetail")
		// if !ok {
		// 	middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
		// 	c.Abort()
		// 	return
		// }
		// serviceDetail := sInterface.(*dao.ServiceDetail)

		appInterface, ok := c.Get("appInfo")
		if !ok {
			middleware.ResponseError(c, 4000, errors.New("AppInfo is not find"))
			c.Abort()
			return
		}
		appInfo := appInterface.(*dao.AppInfo)

		// 统计项
		// 1. 全站统计
		// 2. 服务统计
		// 3. 租户统计
		// 租户统计
		appCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		appCounter.Increase()
		fmt.Println(appInfo.QPD)
		if appInfo.QPD > 0 && appCounter.TotalCount > appInfo.QPD {
			middleware.ResponseError(c, 4002, errors.New(fmt.Sprintf("该租户已限流 limit : %v current : %v",
				appInfo.QPD,
				appCounter.TotalCount)))
			c.Abort()
			return
		}

		// dayappCount, _ := appCounter.GetDayData(time.Now())
		fmt.Printf("{%s serviceCount —— qps:%v total:%v}\n", appInfo.AppID, appCounter.QPS, appCounter.TotalCount)

		c.Next()
	}
}
