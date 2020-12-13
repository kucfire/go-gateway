package realServer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// func main() {
// 	realServer1 := &RealServer{Addr: "127.0.0.1:2003"}
// 	realServer1.Run()
// 	realServer2 := &RealServer{Addr: "127.0.0.1:2004"}
// 	realServer2.Run()

// 	// 监听关闭信号
// 	quit := make(chan os.Signal)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// }
// 	quit := make(chan os.Signal)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// }

// RealServer : real server's Addr
type RealServer struct {
	Addr string
}

// Run : beginning server
func (rs RealServer) Run() {
	log.Println("starting realServer at : ", rs.Addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", rs.HelloFunc)
	mux.HandleFunc("/base/error", rs.ErrorFunc)

	server := &http.Server{
		Addr:         rs.Addr,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

// HelloFunc : func of server
func (rs RealServer) HelloFunc(w http.ResponseWriter, r *http.Request) {
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
func (rs RealServer) ErrorFunc(w http.ResponseWriter, r *http.Request) {
	time.Sleep(6 * time.Second)
	upath := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
}
