package main

import (
	"fmt"
	"go-gateway/load_balance_conf/config"
	"go-gateway/load_balance_conf/proxy"
)

func NewConf() {
	moduleConf, err := config.NewLoadBalanceZkConf("%s",
		"/real_server",
		[]string{"127.0.0.1:2181"},
		map[string]string{},
	)
	if err != nil {
		fmt.Printf("NewLoadBalanceZkConf's error : \n", err)
	}

	loadBalanceObserver := config.NewLoadBalanceObserver(moduleConf)
	moduleConf.Attach(loadBalanceObserver)
	// moduleConf.UpdateConf([]string{"122.11.11"})
	select {}
}

func main() {
	proxy.Run()
	// NewConf()
}
