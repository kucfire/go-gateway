package randomRobin

import (
	"errors"
	"fmt"
	"strings"

	"go-gateway/load_balance_conf/config"
)

// 顺序轮询
type RandomRobinBalance struct {
	curIndex int
	rss      []string
	conf     config.LoadBalanceConf
}

func (r *RandomRobinBalance) Add(key ...string) error {
	if len(key) < 1 {
		return errors.New("key need 1 at least")
	}

	addr := key[0]
	r.rss = append(r.rss, addr)

	return nil
}

func (r *RandomRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RandomRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}

	lens := len(r.rss)

	if r.curIndex >= lens {
		r.curIndex = 0
	}

	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

func (r *RandomRobinBalance) SetConf(conf config.LoadBalanceConf) {
	r.conf = conf
}

func (r *RandomRobinBalance) Update() {
	// 已注册在zk集群上
	if conf, ok := r.conf.(*config.LoadBalanceZkConf); ok {
		fmt.Println("update get conf : ", conf.GetConf())
		r.rss = []string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}
}
