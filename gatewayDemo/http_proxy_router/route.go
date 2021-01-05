package http_proxy_router

import (
	"gatewayDemo/http_proxy_middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {

	// TODO : 优化1 Default会打印对应的请求输出,会消耗一些性能IO，New不会，相应的可以
	// router := gin.Default()
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 使用HTTP中间件
	router.Use(
		http_proxy_middleware.HTTPAccessModeMiddleware(),
		http_proxy_middleware.HTTPFLowCountModeMiddleware(),
		http_proxy_middleware.HTTPFLowLimitModeMiddleware(),
		http_proxy_middleware.HTTPWhiteListModeMiddleware(),
		http_proxy_middleware.HTTPBlackListModeMiddleware(),
		http_proxy_middleware.HTTPHeaderTransferModeMiddleware(),
		http_proxy_middleware.HTTPStripURIModeMiddleware(),
		http_proxy_middleware.HTTPURLRewriteModeMiddleware(),
		http_proxy_middleware.HTTPReverseProxyMiddleware(),
	)
	return router
}
