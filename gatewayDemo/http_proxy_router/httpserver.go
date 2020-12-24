package http_proxy_router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	fmt.Println((&lib.ConfBase.DebugMode == nil))
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}
	go func() {
		log.Printf(" [INFO] http_proxy_run:%s\n", lib.GetStringConf("proxy.http.addr"))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] http_proxy_run:%s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] http_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] http_proxy_stop stopped\n")
}

func HttpsServerRun() {
	fmt.Println((&lib.ConfBase.DebugMode == nil))
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter()
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	go func() {
		log.Printf(" [INFO] https_proxy_run:%s\n", lib.GetStringConf("proxy.https.addr"))
		if err := HttpsSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] https_proxy_run:%s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
		}
	}()
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] https_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] https_proxy_stop stopped\n")
}
