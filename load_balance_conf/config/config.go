package config

import (
	"fmt"
	"go-gateway/load_balance_conf/zook"
)

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	//后期服务发现补充
	Update()
}

type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

type LoadBalanceZkConf struct {
	observers []Observer
	// zk path
	path         string
	zkHosts      []string
	confIpWeight map[string]string
	activeList   []string
	format       string
}

// Attach ： 绑定配置
func (s *LoadBalanceZkConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// NotifyAllObservers : 获取所有监听中的服务器
func (s *LoadBalanceZkConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

// GetConf : 获取ip列表
func (s *LoadBalanceZkConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIpWeight[ip]
		if !ok {
			weight = "50"
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

// 监听
func (s *LoadBalanceZkConf) WatchConf() {
	zkManager := zook.NewZkManager(s.zkHosts)
	zkManager.GetConnect()
	fmt.Println("watchConf")
	chanList, chanErr := zkManager.WatchServerListByPath(s.path)

	go func() {
		defer zkManager.Close()
		for {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changeErr : ", changeErr)
			case changeList := <-chanList:
				fmt.Println("watch node change")
				s.UpdateConf(changeList)
			}
		}
	}()
}

func (s *LoadBalanceZkConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceZkConf(format, path string, zkHosts []string, conf map[string]string) (*LoadBalanceZkConf, error) {
	zkManager := zook.NewZkManager(zkHosts)
	zkManager.GetConnect()
	defer zkManager.Close()

	zlist, err := zkManager.GetServerListByPath(path)
	if err != nil {
		return nil, err
	}
	mconf := &LoadBalanceZkConf{
		path:         path,
		zkHosts:      zkHosts,
		format:       format,
		confIpWeight: conf,
		activeList:   zlist,
	}
	mconf.WatchConf()
	return mconf, nil
}

type Observer interface {
	Update()
}

type LoadBalanceObserver struct {
	ModuleConf *LoadBalanceZkConf
}

func (s *LoadBalanceObserver) Update() {
	fmt.Println("Update get conf : ", s.ModuleConf.GetConf())
}

func NewLoadBalanceObserver(conf *LoadBalanceZkConf) *LoadBalanceObserver {
	return &LoadBalanceObserver{
		ModuleConf: conf,
	}
}
