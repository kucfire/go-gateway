package random

import (
	"errors"
	"fmt"
	"gatewayDemo/reverse_proxy/load_balance_conf/config"
	"math/rand"
	"strings"
)

type RandomBalance struct {
	curIndex int
	rss      []string

	// 后续补充
	conf config.LoadBalanceConf
}

func (r *RandomBalance) Add(key ...string) error {
	if len(key) == 0 {
		return errors.New("key need 1 at least")
	}

	addr := key[0]
	r.rss = append(r.rss, addr)

	return nil
}

func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) SetConf(conf config.LoadBalanceConf) {
	r.conf = conf
}

func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	// 以r.rss的长度为随即范围
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

func (r *RandomBalance) Update() {
	// 更新
	if conf, ok := r.conf.(*config.LoadBalanceZkConf); ok {
		fmt.Println("RandomBalance get conf : ", conf.GetConf())
		r.rss = []string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}

	if conf, ok := r.conf.(*config.LoadBalanceZkCheckConf); ok {
		fmt.Println("RandomBalance get conf : ", conf.GetConf())
		r.rss = []string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}
}
