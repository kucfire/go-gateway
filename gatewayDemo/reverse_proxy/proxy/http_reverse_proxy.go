package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gatewayDemo/middleware"
	"gatewayDemo/reverse_proxy/load_balance_conf/config"

	"github.com/gin-gonic/gin"
)

// var (
// 	addr      = "127.0.0.1:2002"
// 	transport = &http.Transport{
// 		DialContext: (&net.Dialer{
// 			Timeout:   30 * time.Second, //连接超时
// 			KeepAlive: 30 * time.Second, //长连接超时时间
// 		}).DialContext,
// 		MaxIdleConns:          100,              //最大空闲连接
// 		IdleConnTimeout:       90 * time.Second, //空闲超时时间
// 		TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
// 		ExpectContinueTimeout: 1 * time.Second,  //100-continue超时时间
// 	}
// )

func NewLoadBalanceReverseProxy(c *gin.Context, lb config.LoadBalance, transport *http.Transport) *httputil.ReverseProxy {
	//请求协调者
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.URL.String())
		// TODO ： 优化点3 注释掉对应的打印操作，节省io开销
		// fmt.Println(req.URL.String())
		// fmt.Println(nextAddr)
		if err != nil {
			panic("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			panic(err)
		}
		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		// fmt.Println(target.Scheme)
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		// fmt.Println(target.Host)
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

		// TODO ： 优化点2 压测时不需要进行更改操作
		// var payload []byte
		// payload, err := ioutil.ReadAll(resp.Body)
		// var readErr error

		//todo 部分章节功能补充3
		//todo 兼容gzip压缩
		// if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		// 	gr, err := gzip.NewReader(resp.Body)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	payload, readErr = ioutil.ReadAll(gr)
		// 	resp.Header.Del("Content-Encoding")
		// } else {
		// 	payload, readErr = ioutil.ReadAll(resp.Body)
		// }
		// if readErr != nil {
		// 	return readErr
		// }

		// //异常请求时设置StatusCode
		// if resp.StatusCode != 200 {
		// 	payload = []byte("StatusCode error:" + string(payload))
		// }

		//todo 部分章节功能补充4
		//todo 因为预读了数据所以内容重新回写
		// c.Set("status_code", resp.StatusCode)
		// c.Set("payload", payload)
		// resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		// resp.ContentLength = int64(len(payload))
		// resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}

	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		//todo record error log
		middleware.ResponseError(c, 999, err)
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
