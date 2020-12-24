package main

import (
	"flag"
	"gatewayDemo/http_proxy_router"
	"gatewayDemo/router"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/e421083458/golang_common/lib"
)

// endpoint dashboard后台管理 proxy_server代理服务器
// config 。/conf/prod/ 对应配置文件夹
var (
	endpoint = flag.String(
		"endpoint",                           // flag name 标志命名
		"",                                   // flag value 值
		"input endpoint dashboard or server", // flag usge 用法
	)
	conf = flag.String(
		"conf",                               // flag name 标志命名
		"",                                   // flag value 值
		"input config file like ./conf/dev/", // flag usge 用法
	)
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *conf == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 如果是跑在windows os上，则需要写死路径
	// 在cmd跟powershell中没办法输入带.的参数
	if runtime.GOOS == "windows" {
		*conf = "./conf/dev/"
	}

	if *endpoint == "dashboard" {
		// 如果configPath为空，则从命令行中‘-config=。/conf/prod/‘中读取。
		lib.InitModule(*conf, []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		lib.InitModule(*conf, []string{"base", "mysql", "redis"})
		defer lib.Destroy()

		go func() {
			http_proxy_router.HttpServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		http_proxy_router.HttpServerStop()
	}
}

func main2() {
	lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
	defer lib.Destroy()
	http_proxy_router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	http_proxy_router.HttpServerStop()
}
