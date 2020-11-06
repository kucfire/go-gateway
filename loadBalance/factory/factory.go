package factory

import (
	"go-gateway/loadBalance/config"
	"go-gateway/loadBalance/demo/hashrandom"
	"go-gateway/loadBalance/demo/random"
	"go-gateway/loadBalance/demo/randomRobin"
	"go-gateway/loadBalance/demo/weightroundrobin"
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
