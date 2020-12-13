package main

import (
	"fmt"
	"go-gateway/middlewareDemo/middleware"
	"go-gateway/middlewareDemo/proxy"
	"go-gateway/middlewareDemo/realServer"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

// 监听127.0.0.1:2002的地址（ip+port）
var addr = "127.0.0.1:2002"

func middleWareAssembly() {
	//
	reverseProxy := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003"
		url1, err1 := url.Parse(rs1)
		if err1 != nil {
			panic(err1)
		}

		rs2 := "http://127.0.0.1:2004"
		url2, err2 := url.Parse(rs2)
		if err2 != nil {
			panic(err2)
		}

		urls := []*url.URL{url1, url2}

		return proxy.NewMultipleHostsReverseProxy(c, urls)
	}
	log.Println("starting httpserver at : ", addr)

	// 初始化方法数组路由器
	sliceRouter := middleware.NewSliceRouter()

	// 中间件可充当业务逻辑代码
	sliceRouter.Group("/base").Use(middleware.TraceLogSliceMW(),
		func(c *middleware.SliceRouterContext) {
			c.RW.Write([]byte("test func"))
		})
	routerHandler := middleware.NewSliceRouterHandler(nil, sliceRouter)

	// 请求到反向代理
	sliceRouter.Group("/").Use(middleware.TraceLogSliceMW(),
		func(c *middleware.SliceRouterContext) {
			fmt.Println("reverseProxy")
			reverseProxy(c).ServeHTTP(c.RW, c.Req)
		})

	routerHandler = middleware.NewSliceRouterHandler(nil, sliceRouter)
	log.Fatal(http.ListenAndServe(addr, routerHandler))
}

func realserver() {
	realServer1 := &realServer.RealServer{Addr: "127.0.0.1:2003"}
	realServer1.Run()
	realServer2 := &realServer.RealServer{Addr: "127.0.0.1:2004"}
	realServer2.Run()

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func middlewareFusing() {
	coreFunc := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003"
		url1, err1 := url.Parse(rs1)
		if err1 != nil {
			panic(err1)
		}

		rs2 := "http://127.0.0.1:2004"
		url2, err2 := url.Parse(rs2)
		if err2 != nil {
			panic(err2)
		}

		urls := []*url.URL{url1, url2}

		return proxy.NewMultipleHostsReverseProxy(c, urls)
	}
	log.Println("starting httpserver at : ", addr)

	// 初始化方法数组路由器
	sliceRouter := middleware.NewSliceRouter()
	sliceRouter.Group("/base").Use(middleware.RateLimiter())
	routerHandler := middleware.NewSliceRouterHandler(coreFunc, sliceRouter)
	log.Fatal(http.ListenAndServe(addr, routerHandler))
}

func main() {
	// 中间件demo
	// go middleWareAssembly()
	// 熔断中间件方案
	go middlewareFusing()

	realserver()

}
