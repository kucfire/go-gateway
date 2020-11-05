package realServer

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type RealServer struct {
	Addr string
}

// Set : set server addr
func (rs *RealServer) Set(addr string) {
	rs.Addr = addr
}

func (rs *RealServer) Run() {
	log.Println("starting httpServer at : ", rs.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", rs.helloHandle)
}

func (rs *RealServer) helloHandle(w http.ResponseWriter, r *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", rs.Addr)
	io.WriteString(w, upath)
}

func (rs *RealServer) ErrorHandle(w http.ResponseWriter, r *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}
