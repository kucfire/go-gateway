package config

import "testing"

func TestLoadBalanceObserver(t *testing.T) {
	moduleConf, err := NewLoadBalanceZkConf("%s",
		"/real_server",
		[]string{"127.0.0.1:2181"},
		map[string]string{"127.0.0.1:2003": "20"},
	)
	if err != nil {
		t.Error("NewLoadBalanceZkConf's error : \n", err)
	}

	loadBalanceObserver := NewLoadBalanceObserver(moduleConf)
	moduleConf.Attach(loadBalanceObserver)
	moduleConf.UpdateConf([]string{"122.11.11"})
	select {}
}
