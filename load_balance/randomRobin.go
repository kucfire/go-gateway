package load_balance

import (
	"errors"
)

type RandomRobinBalance struct {
	curIndex int
	rss      []string
	conf     LoadBalanceConf
}

// Add : add addr in RandomBalance.rss
func (r *RandomRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}

	addr := params[0]
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

	// 如果curIndex的值大于等于服务器连接池总数的话会从零开始
	if r.curIndex >= lens {
		r.curIndex = 0
	}

	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens // double setting如果curIndex的值大于等于服务器连接池总数的话会从零开始
	return curAddr
}
