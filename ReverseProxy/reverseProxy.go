package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var addr = "127.0.0.1:2002"

func main() {
	// 127.0.0.1:2002/xxx
	// actual 127.0.0.1:2003/base/xxx
	rs1 := "http://127.0.0.1:2003/base"
	url1, err := url.Parse(rs1)
	if err != nil {
		log.Println("err:", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url1)
	log.Println("Starting httpServer at " + addr)
	log.Fatal(http.ListenAndServe(":2002", proxy))

}
