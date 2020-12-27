package factory

import (
	"gatewayDemo/reverse_proxy/load_balance_conf/config"
	"gatewayDemo/reverse_proxy/load_balance_conf/load_balance/demo/hashrandom"
	"gatewayDemo/reverse_proxy/load_balance_conf/load_balance/demo/random"
	"gatewayDemo/reverse_proxy/load_balance_conf/load_balance/demo/randomRobin"
	"gatewayDemo/reverse_proxy/load_balance_conf/load_balance/demo/weightroundrobin"
)

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

func LoadBalanceFactoryWithConf(lbtype LbType, mconf config.LoadBalanceConf) config.LoadBalance {
	switch lbtype {
	case LbConsistentHash:
		lb := hashrandom.NewConsistentHashmapBalance(10, nil)
		lb.SetConf(mconf)
		mconf.Attach(lb)
		return lb
	case LbRoundRobin:
		lb := &randomRobin.RandomRobinBalance{}
		lb.SetConf(mconf)
		mconf.Attach(lb)
		return lb
	case LbWeightRoundRobin:
		lb := &weightroundrobin.WeightRoundRobinBalance{}
		lb.SetConf(mconf)
		mconf.Attach(lb)
		return lb
	case LbRandom:
		lb := &random.RandomBalance{}
		lb.SetConf(mconf)
		mconf.Attach(lb)
		return lb
	default:
		lb := &random.RandomBalance{}
		lb.SetConf(mconf)
		mconf.Attach(lb)
		return lb
	}
}
