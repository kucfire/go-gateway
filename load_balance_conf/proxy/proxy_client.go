package proxy

import (
	"go-gateway/load_balance_conf/config"
	"go-gateway/load_balance_conf/load_balance/factory"
	"go-gateway/middlewareDemo/middleware"
	"log"
	"net/http"
)

func RunClient() {
	// zk设置
	mconf, err := config.NewLoadBalanceZkCheckConf("http://%s/base",
		map[string]string{"127.0.0.1:2003": "20", "127.0.0.1:2004": "20"})
	if err != nil {
		panic(err)
	}

	// 负载均衡设置
	rb := factory.LoadBalanceFactoryWithConf(factory.LbWeightRoundRobin, mconf)

	// 代理
	proxy := NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)
	log.Println("starting proxy httpserver : ", addr)

	// 监听服务器
	log.Fatal(http.ListenAndServe(addr, proxy))
}
