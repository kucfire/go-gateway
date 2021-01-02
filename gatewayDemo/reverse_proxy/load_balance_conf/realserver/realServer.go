package main

import (
	"fmt"
	"gatewayDemo/reverse_proxy/load_balance_conf/zook"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	realServer1 := &RealServer{Addr: "127.0.0.1:2003"}
	realServer1.Run()
	time.Sleep(2 * time.Second)
	realServer2 := &RealServer{Addr: "127.0.0.1:2004"}
	realServer2.Run()
	time.Sleep(2 * time.Second)

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// RealServer : real server's Addr
type RealServer struct {
	Addr string
}

// Run : beginning server
func (rs *RealServer) Run() {
	log.Println("starting realServer at : ", rs.Addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", rs.HelloFunc)
	mux.HandleFunc("/base/error", rs.ErrorFunc)

	// mux.HandleFunc("/test_http_string/aaa", rs.TimeoutErrorFunc)

	server := &http.Server{
		Addr:         rs.Addr,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}
	go func() {
		// 链接zookeeper
		zkManager := zook.NewZkManager([]string{"127.0.0.1:2181"})
		err := zkManager.GetConnect()
		if err != nil {
			fmt.Println("connect zk error : ", err)
		}
		defer zkManager.Close()

		//
		err = zkManager.RegistServerPath("/real_server", rs.Addr)
		if err != nil {
			fmt.Println("regist node error : ", err)
		}

		zlist, err := zkManager.GetServerListByPath("/real_server")
		fmt.Println(zlist)

		// 监听服务器
		log.Fatal(server.ListenAndServe())
	}()
}

// HelloFunc : func of server
func (rs *RealServer) HelloFunc(w http.ResponseWriter, r *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", rs.Addr, r.URL.Path)
	realIP := fmt.Sprintf(
		"RemoteAddr=%s,X-Forwarded-For=%v,x-Real-Ip=%v\n",
		r.RemoteAddr,
		r.Header.Get("X-Forwarded-For"),
		r.Header.Get("X-Real-Ip"),
	)
	header := fmt.Sprintf("headers =%v\n", r.Header)
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
	io.WriteString(w, header)
}

// ErrorFunc : error handle
func (rs *RealServer) ErrorFunc(w http.ResponseWriter, r *http.Request) {
	time.Sleep(6 * time.Second)
	upath := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
}

// TimeoutErrorFunc : timeout error handle
func (rs *RealServer) TimeoutErrorFunc(w http.ResponseWriter, r *http.Request) {
	time.Sleep(6 * time.Second)
	upath := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
}
