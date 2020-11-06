package realServer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// RealServer : Body
type RealServer struct {
	Addr string
}

// Set : set server addr
func (rs *RealServer) Set(addr string) {
	rs.Addr = addr
}

// Run : start server
func (rs *RealServer) Run() {
	log.Println("starting httpServer at : ", rs.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", rs.helloHandle)
	mux.HandleFunc("/base/error", rs.errorHandle)

	server := &http.Server{
		Addr:         rs.Addr,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

func (rs *RealServer) helloHandle(w http.ResponseWriter, r *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", rs.Addr, r.URL.Path)
	io.WriteString(w, upath)
}

func (rs *RealServer) errorHandle(w http.ResponseWriter, r *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", rs.Addr, r.URL.Path) + "error handler\n"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}
