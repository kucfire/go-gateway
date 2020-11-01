package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

var (
	port       = "2002"
	proxy_addr = "http://127.0.0.1:2003"
)

func handler(w http.ResponseWriter, req *http.Request) {
	// Step1 解析代理地址，并更改请求体的协议和主机
	proxy, err := url.Parse(proxy_addr)
	if err != nil {
		log.Print(err)
	}
	req.URL.Scheme = proxy.Scheme
	req.URL.Host = proxy.Host

	// Step2 请求下游
	transprot := http.DefaultTransport
	resp, err := transprot.RoundTrip(req)
	if err != nil {
		log.Print(err)
		return
	}

	// Step3 把下游请求内容返回给上游
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	defer resp.Body.Close()
	bufio.NewReader(resp.Body).WriteTo(w)
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("Start serving on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
