package http_proxy_middleware

import (
	"gatewayDemo/dao"
	"gatewayDemo/middleware"
	"gatewayDemo/public"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// jwt auth token
func HTTPJwtOauthTokenModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sInterface, ok := c.Get("serviceDetail")
		if !ok {
			middleware.ResponseError(c, 1002, errors.New("serviceDetail is not find"))
			c.Abort()
			return
		}
		serviceDetail := sInterface.(*dao.ServiceDetail)

		// decode jwt token
		// app_id 与app_list 取得 appInfo
		// appInfo 放到 gin.context
		// token := strings.Split(c.GetHeader("Authorization"), " ")[1]
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		matched := false
		if token != "" {
			claim, err := public.JwtDecode(token)
			if err != nil {
				middleware.ResponseError(c, 1003, err)
				c.Abort()
				return
			}
			// fmt.Println("claims.Issuer : ", claim.Issuer)
			appList := dao.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claim.Issuer {
					c.Set("appInfo", appInfo)
					matched = true
					break
				}
			}
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && !matched {
			middleware.ResponseError(c, 1004, errors.New("not mactched vail app"))
			c.Abort()
			return
		}
		c.Next()
	}
}
