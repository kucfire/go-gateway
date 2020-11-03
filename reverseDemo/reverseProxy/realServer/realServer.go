package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// real server setting
	rs1 := &RealServer{
		Addr: "127.0.0.1:2003",
	}
	rs1.Run()
	rs2 := &RealServer{
		Addr: "127.0.0.1:2004",
	}
	rs2.Run()

	// listen close signal
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

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
