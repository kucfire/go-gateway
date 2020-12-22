package factory

import (
	"go-gateway/load_balance_conf/config"
	"go-gateway/load_balance_conf/load_balance/demo/hashrandom"
	"go-gateway/load_balance_conf/load_balance/demo/random"
	"go-gateway/load_balance_conf/load_balance/demo/randomRobin"
	"go-gateway/load_balance_conf/load_balance/demo/weightroundrobin"
)

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

func LoadBalanceFactory(lbtype LbType) config.LoadBalance {
	switch lbtype {
	case LbConsistentHash:
		return hashrandom.NewConsistentHashmapBalance(10, nil)
	case LbRoundRobin:
		return &randomRobin.RandomRobinBalance{}
	case LbWeightRoundRobin:
		return &weightroundrobin.WeightRoundRobinBalance{}
	case LbRandom:
		return &random.RandomBalance{}
	default:
		return &random.RandomBalance{}
	}
}
