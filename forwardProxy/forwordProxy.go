package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct{}

// ServeHTTP : Handle Body, receive
func (p *Pxy) ServeHTTP(r http.ResponseWriter, q *http.Request) {
	fmt.Printf("Receved request %s %s %s\n", q.Method, q.Host, q.RemoteAddr)
	transport := http.DefaultTransport

	// step1 : 浅拷贝对象，然后就再新增属性数据
	outReq := new(http.Request)
	*outReq = *q
	// split Host and Port
	if clientIP, _, err := net.SplitHostPort(q.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forward-For"]; ok {
			clientIP = strings.Join(prior, ",") + "," + clientIP
		}
		fmt.Println(clientIP)
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// step2 : 请求下游
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		r.WriteHeader(http.StatusBadGateway)
		return
	}

	// step3 : 把下游请求内容返回给上游
	for key, value := range res.Header {
		for _, v := range value {
			r.Header().Add(key, v)
		}
	}
	r.WriteHeader(res.StatusCode)
	io.Copy(r, res.Body)
	res.Body.Close()
}

func main() {
	fmt.Println("serve on :8080")
	http.Handle("/", &Pxy{})
	http.ListenAndServe("0.0.0.0:8080", nil)
}
