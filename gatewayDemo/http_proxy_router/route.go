package http_proxy_router

import (
	"gatewayDemo/controller"
	"gatewayDemo/http_proxy_middleware"
	"gatewayDemo/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {

	// TODO : 优化1 Default会打印对应的请求输出,会消耗一些性能io，New不会，相应的可以
	router := gin.Default()
	// router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	oauthRouter := router.Group("/oauth")
	oauthRouter.Use(middleware.TranslationMiddleware())
	{
		controller.OauthRegister(oauthRouter)
	}

	// 使用HTTP中间件
	// root := router.Group("/")
	router.Use(
		//
		http_proxy_middleware.HTTPAccessModeMiddleware(),
		//
		http_proxy_middleware.HTTPFLowCountModeMiddleware(),
		http_proxy_middleware.HTTPFLowLimitModeMiddleware(),
		// 权限校验
		http_proxy_middleware.HTTPJwtOauthTokenModeMiddleware(),
		http_proxy_middleware.HTTPJwtFLowCountModeMiddleware(),
		http_proxy_middleware.HTTPFJwtLowLimitModeMiddleware(),
		http_proxy_middleware.HTTPWhiteListModeMiddleware(),
		http_proxy_middleware.HTTPBlackListModeMiddleware(),
		//
		http_proxy_middleware.HTTPHeaderTransferModeMiddleware(),
		http_proxy_middleware.HTTPStripURIModeMiddleware(),
		http_proxy_middleware.HTTPURLRewriteModeMiddleware(),
		// 代理层
		http_proxy_middleware.HTTPReverseProxyMiddleware(),
	)

	return router
}
