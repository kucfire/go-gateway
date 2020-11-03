package load_balance

import "errors"

type RandomBalance struct {
	curIndex int
	rss      []string
	// conf     LoadBalanceConf
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	return nil
}
