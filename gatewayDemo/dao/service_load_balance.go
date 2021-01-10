package dao

import (
	"fmt"
	"gatewayDemo/public"
	"gatewayDemo/reverse_proxy/load_balance_conf/config"
	"gatewayDemo/reverse_proxy/load_balance_conf/load_balance/factory"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

func init() {
	LoadBalanceHandler = NewLoadBalancer()
	TransportorHandler = NewTransportor()
}

type ServiceLoadBalance struct {
	ID                     int64  `json:"id" gorm:"primary_key"`
	ServiceID              int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	CheckMethod            int    `json:"check_method" gorm:"column:check_method" description:"检查方法 0=tcpchk，检测端口是否握手成功"`
	CheckTimeout           int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间，单位s"`
	CheckInterval          int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔，单位s"`
	RoundType              int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash"`
	IPList                 string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList             string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList             string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"建立连接超时，单位s"`
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"获取header超时，单位s"`
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"链接最大空闲时间，单位s"`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"最大空闲链接数"`
}

func (t *ServiceLoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *ServiceLoadBalance) Find(c *gin.Context, tx *gorm.DB, search *ServiceLoadBalance) (*ServiceLoadBalance, error) {
	out := &ServiceLoadBalance{}
	// 将查询出来的结果放入out结构提里面去
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (t *ServiceLoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	// 将ad保存进数据库
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}

func (t *ServiceLoadBalance) GetIPList() []string {
	return strings.Split(t.IPList, ",")
}

func (t *ServiceLoadBalance) GetWeightList() []string {
	return strings.Split(t.WeightList, ",")
}

var LoadBalanceHandler *LoadBalancer

type LoadBalanceItem struct {
	LoadBalance config.LoadBalance
	ServiceName string
}

type LoadBalancer struct {
	// 服务较多时，可以用
	LoadBalanceMap   map[string]*LoadBalanceItem
	LoadBalanceSlice []*LoadBalanceItem
	Locker           sync.RWMutex
	// init             sync.Once
	// errMsg           error
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBalanceMap:   map[string]*LoadBalanceItem{},
		LoadBalanceSlice: []*LoadBalanceItem{},
		Locker:           sync.RWMutex{},
	}
}

func (lbr *LoadBalancer) GetLoadBalance(service *ServiceDetail) (config.LoadBalance, error) {
	// 遍历Slice看是否已存在相同服务名称的负载均衡设置
	for _, loadBalanceItem := range lbr.LoadBalanceSlice {
		if loadBalanceItem.ServiceName == service.Info.ServiceName {
			return loadBalanceItem.LoadBalance, nil
		}
	}

	// 判断协议是否需要添加加密协议
	schema := "http://"
	if service.HTTPRule != nil && service.HTTPRule.NeedHTTPS == 1 {
		schema = "https://"
	}
	if service.Info.LoadType == public.LoadTypeTCP || service.Info.LoadType == public.LoadTypeGRPC {
		schema = ""
	}

	// prefix := ""
	// if service.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
	// 	prefix = service.HTTPRule.Rule
	// }

	// ip列表设置
	ipList := service.LoadBalance.GetIPList()
	weightList := service.LoadBalance.GetWeightList()
	ipconf := map[string]string{}
	for i, list := range ipList {
		ipconf[list] = weightList[i]
	}
	// fmt.Println(ipconf)

	// zk设置
	mconf, err := config.NewLoadBalanceZkCheckConf( // "http://%s/base",
		fmt.Sprintf("%s%s", schema, "%s"), ipconf)
	if err != nil {
		return nil, err
	}

	// 负载均衡设置
	// lb := factory.LoadBalanceFactoryWithConf(factory.LbWeightRoundRobin, mconf)
	lb := factory.LoadBalanceFactoryWithConf(factory.LbType(service.LoadBalance.RoundType), mconf)
	item := &LoadBalanceItem{
		LoadBalance: lb,
		ServiceName: service.Info.ServiceName,
	}

	// save to map and slice
	lbr.LoadBalanceSlice = append(lbr.LoadBalanceSlice, item)
	lbr.Locker.Lock()
	defer lbr.Locker.Unlock()
	lbr.LoadBalanceMap[service.Info.ServiceName] = item
	return lb, nil
}

// 连接池设置
var TransportorHandler *Transportor

type Transportor struct {
	// 服务较多时，可以用
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Locker         sync.RWMutex
	// init             sync.Once
	// errMsg           error
}

type TransportItem struct {
	Trans       *http.Transport
	ServiceName string
}

func NewTransportor() *Transportor {
	return &Transportor{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Locker:         sync.RWMutex{},
	}
}

func (t *Transportor) GetTransportor(service *ServiceDetail) (*http.Transport, error) {
	// 遍历Slice看是否已存在相关的transport设置
	for _, transItem := range t.TransportSlice {
		if transItem.ServiceName == service.Info.ServiceName {
			return transItem.Trans, nil
		}
	}

	// TODO : 优化点五
	if service.LoadBalance.UpstreamHeaderTimeout == 0 {
		service.LoadBalance.UpstreamHeaderTimeout = 30
	}
	if service.LoadBalance.UpstreamIdleTimeout == 0 {
		service.LoadBalance.UpstreamIdleTimeout = 90
	}
	if service.LoadBalance.UpstreamConnectTimeout == 0 {
		service.LoadBalance.UpstreamConnectTimeout = 30
	}
	if service.LoadBalance.UpstreamMaxIdle == 0 {
		service.LoadBalance.UpstreamMaxIdle = 100
	}

	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(service.LoadBalance.UpstreamConnectTimeout) * time.Second, //连接超时
			KeepAlive: 30 * time.Second,                                                        //长连接超时时间
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          service.LoadBalance.UpstreamMaxIdle,                                  //最大空闲连接
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout) * time.Second, //空闲超时时间
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout) * time.Second,
	}
	transItem := &TransportItem{
		Trans:       trans,
		ServiceName: service.Info.ServiceName,
	}

	// save to map and slice
	t.TransportSlice = append(t.TransportSlice, transItem)
	// 服务多的话使用map会快一点
	t.Locker.Lock()
	defer t.Locker.Unlock()
	t.TransportMap[service.Info.ServiceName] = transItem

	// 负载均衡设置
	return trans, nil
}
