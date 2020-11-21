package dao

import (
	"gatewayDemo/public"
	"strings"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceLoadBalance struct {
	Id                     int64  `json:"id" gorm:"primary_key"`
	ServiceID              int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	CheckMethod            int    `json:"check_method" gorm:"column:check_method" description:"检查方法 0=tcpchk，检测端口是否握手成功"`
	CheckTimeout           int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间，单位s"`
	CheckInterval          int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔，单位s"`
	RoundType              int8   `json:"round_type" gorm:"column:round_type" description:"轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash"`
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
