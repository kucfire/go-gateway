package config

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

// type LoadBalanceConf interface {
// 	Attach(o Observer)
// 	GetConf() []string
// 	WatchConf()
// 	UpdateConf(conf []string)
// }

const (
	//default check setting
	//默认查看配置
	DefaultCheckMethod    = 0 //
	DefaultCheckTimeout   = 2 // 超时时间，单位s
	DefaultCheckMaxErrNum = 2 // 最大容错值
	DefaultCheckInterval  = 5 // 间隔，单位s
)

type LoadBalanceZkCheckConf struct {
	observers []Observer
	// zk path
	confIpWeight map[string]string
	activeList   []string
	format       string
}

// Attach ： 绑定配置
func (s *LoadBalanceZkCheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// NotifyAllObservers : 获取所有监听中的服务器
func (s *LoadBalanceZkCheckConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

// GetConf : 获取ip列表
func (s *LoadBalanceZkCheckConf) GetConf() []string {
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
func (s *LoadBalanceZkCheckConf) WatchConf() {
	fmt.Println("watchConf")
	go func() {
		// 计数：错误次数
		confIpErrNum := map[string]int{}
		for {
			changeList := []string{}
			for item, _ := range s.confIpWeight {
				// tcp连接测试
				conn, err := net.DialTimeout("tcp", item, time.Duration(DefaultCheckTimeout)*time.Second)
				// TODO http连接测试
				//
				if err == nil {
					// 测试成功，连接可用
					conn.Close()
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] = 0
					}
				}
				if err != nil {
					// 累加失败次数
					confIpErrNum[item]++
				}
				// 当错误次数不超过默认最大容错值，则加入list里面
				if confIpErrNum[item] < DefaultCheckMaxErrNum {
					changeList = append(changeList, item)
				}
			}

			// 将元素以递增方式排列，我也不知道为啥这里需要sort一下
			// 可能是想讲ip排列后方便查看以及排查确实的port或者ip
			// 可用作对照
			sort.Strings(changeList)
			sort.Strings(s.activeList)
			if !reflect.DeepEqual(changeList, s.activeList) {
				fmt.Println(changeList)
				fmt.Println(s.activeList)
				s.UpdateConf(changeList)
			}
			// 心跳检测间隔
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

func (s *LoadBalanceZkCheckConf) UpdateConf(conf []string) {
	fmt.Println("UpdateConf", conf)
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceZkCheckConf(format string, conf map[string]string) (*LoadBalanceZkCheckConf, error) {
	alist := []string{}

	// 默认初始化
	for item, _ := range conf {
		alist = append(alist, item)
	}

	mconf := &LoadBalanceZkCheckConf{
		format:       format,
		activeList:   alist,
		confIpWeight: conf,
	}

	// 监听
	mconf.WatchConf()
	return mconf, nil
}
