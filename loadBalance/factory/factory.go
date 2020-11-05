package factory

import "go-gateway/loadBalance/config"

type LbType int

const (
	LbRandom = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

func LoadBalanceFactory(lbtype LbType) config.LoadBalance {
	switch lbtype {
	default:
		return &
	}
}
