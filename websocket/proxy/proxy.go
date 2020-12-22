package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"go-gateway/loadBalance/config"
	"go-gateway/loadBalance/factory"
	"go-gateway/middlewareDemo/middleware"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	addr      = "127.0.0.1:2002"
	transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, //连接超时
			KeepAlive: 30 * time.Second, //长连接超时时间
		}).DialContext,
		MaxIdleConns:          100,              //最大空闲连接
		IdleConnTimeout:       90 * time.Second, //空闲超时时间
		TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  //100-continue超时时间
	}
)

func NewMultipleHostsReverseProxy(targets config.LoadBalance) *httputil.ReverseProxy {
	// 请求协议者
	director := func(req *http.Request) {

		nextAddr, errGet := targets.Get(req.RemoteAddr)
		if errGet != nil {
			log.Fatal("get next addr fail" + errGet.Error())
		}

		target, errParse := url.Parse(nextAddr)
		if errParse != nil {
			fmt.Println(nextAddr)
			log.Fatal(errParse)
		}

		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		// url地址重写：重写前：/aa 重写后：/base/aa
		req.URL.Path = singleJoiningSlash(target.Path, ""+req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
		//只在第一代理中设置此header头
	}

	// 更改内容
	modifyFunc := func(resp *http.Response) error {
		// 请求一下命令 : curl "http://127.0.0.1:2002/error"
		// todo 部分章节功能补充2
		// todo 兼容websocket
		if strings.Contains(resp.Header.Get("connection"), "Upgrade") {
			return nil
		}
		var payload []byte
		var readErr error

		// todo 部分章节功能补充3
		// todo 兼容gzip压缩
		if strings.Contains(resp.Header.Get("content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			payload, readErr = ioutil.ReadAll(gr)
			resp.Header.Del("content-Encoding")
		} else {
			payload, readErr = ioutil.ReadAll(resp.Body)
		}
		if readErr != nil {
			return readErr
		}

		// 异常请求时设置StatusCode
		if resp.StatusCode != 200 {
			payload = []byte("StatusCode error:" + string(payload))
		}

		// todo 部分章节功能补充4
		// todo 因为预读了数据所以内容重新回写
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		resp.ContentLength = int64(len(payload))
		resp.Header.Set("content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}

	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "ErrorHandler error : "+err.Error(), 500)
	}

	return &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}

func NewLoadBalanceReverseProxy(c *middleware.SliceRouterContext, lb config.LoadBalance) *httputil.ReverseProxy {
	//请求协调者
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.URL.String())
		if err != nil {
			log.Fatal("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			log.Fatal(err)
		}
		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		fmt.Println(target.Host)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	//更改内容
	modifyFunc := func(resp *http.Response) error {
		//todo 部分章节功能补充2
		//todo 兼容websocket
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		var payload []byte
		// payload, err := ioutil.ReadAll(resp.Body)
		var readErr error

		//todo 部分章节功能补充3
		//todo 兼容gzip压缩
		if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			payload, readErr = ioutil.ReadAll(gr)
			resp.Header.Del("Content-Encoding")
		} else {
			payload, readErr = ioutil.ReadAll(resp.Body)
		}
		if readErr != nil {
			return readErr
		}

		//异常请求时设置StatusCode
		if resp.StatusCode != 200 {
			payload = []byte("StatusCode error:" + string(payload))
		}

		//todo 部分章节功能补充4
		//todo 因为预读了数据所以内容重新回写
		c.Set("status_code", resp.StatusCode)
		c.Set("payload", payload)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		resp.ContentLength = int64(len(payload))
		resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}

	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		//todo record error log
		fmt.Println(err)
	}

	return &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// func Run2() {
// 	rb := factory.LoadBalanceFactory(factory.LbWeightRoundRobin)
// 	rb.Add("http://127.0.0.1:2003", "50")
// 	proxy := NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)
// 	// proxy := NewLoadBalanceReverseProxy(newSliceRouterContext(), rb)
// 	log.Println("Starting httpserver at " + addr)
// 	log.Fatal(http.ListenAndServe(addr, proxy))
// }

func Run() {
	coreFunc := func(c *middleware.SliceRouterContext) http.Handler {
		rb := factory.LoadBalanceFactory(factory.LbWeightRoundRobin)
		rb.Add("http://127.0.0.1:2003", "50")
		rb.Add("http://127.0.0.1:2004", "60")

		return NewLoadBalanceReverseProxy(c, rb)
		// proxy := NewMultipleHostsReverseProxy(rb)
	}

	// 初始化方法数组路由器
	sliceRouter := middleware.NewSliceRouter()
	// sliceRouter.Group("/").Use(middleware.RateLimiter())

	routerHandler := middleware.NewSliceRouterHandler(coreFunc, sliceRouter)
	log.Println("Start http server at:" + addr)
	log.Fatal(http.ListenAndServe(addr, routerHandler))
}