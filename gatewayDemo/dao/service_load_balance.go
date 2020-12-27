package dao

import (
	"fmt"
	"gatewayDemo/public"
	"gatewayDemo/reverse_proxy/load_balance_conf/config"
	"gatewayDemo/reverse_proxy/load_balance_conf/load_balance/factory"
	"strings"
	"sync"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

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

func init() {
	LoadBalanceHandler = NewLoadBalancer()
}

type LoadBalancer struct {
	// 服务较多时，可以用
	LoadBalanceMap   map[string]config.LoadBalance
	LoadBalanceSlice []config.LoadBalance
	Locker           sync.RWMutex
	// init             sync.Once
	// errMsg           error
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBalanceMap:   map[string]config.LoadBalance{},
		LoadBalanceSlice: []config.LoadBalance{},
		Locker:           sync.RWMutex{},
	}
}

func (lbr *LoadBalancer) GetLoadBalance(service *ServiceDetail) (config.LoadBalance, error) {
	// 判断协议是否需要添加加密协议
	schema := "http"
	if service.HTTPRule.NeedHTTPS == 1 {
		schema = "https"
	}

	prefix := ""
	if service.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
		prefix = service.HTTPRule.Rule
	}

	// ip列表设置
	ipList := service.LoadBalance.GetIPList()
	weightList := service.LoadBalance.GetWeightList()
	ipconf := map[string]string{}
	for i, list := range ipList {
		ipconf[list] = weightList[i]
	}

	// zk设置
	mconf, err := config.NewLoadBalanceZkCheckConf( // "http://%s/base",
		fmt.Sprintf("%s://%s%s", schema, prefix), ipconf)
	if err != nil {
		return nil, err
	}

	// 负载均衡设置
	return factory.LoadBalanceFactoryWithConf(factory.LbWeightRoundRobin, mconf), nil
}

// 连接池设置
