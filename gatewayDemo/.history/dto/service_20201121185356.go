package dto

import (
	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
)

// ServiceListInput : 服务列表输入结构体
type ServiceListInput struct {
	// 关键词
	Info string `json:"info" form:"info" comment:"关键词" example:"" validate:""`
	// 页数
	PageNo int `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`
	// 每页条数
	PageSize int `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"`
}

func (params *ServiceListInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ServiceListOutput ： 服务列表输出结构体
type ServiceListOutput struct {
	// 总数
	Total int64 `json:"total" form:"total" comment:"总数" example:"" validate:""`
	// 列表
	List []ServiceListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`
}

// ServiceListItemOutput ： 服务列表中的表单输出结构体
type ServiceListItemOutput struct {
	// 总数
	ID int64 `json:"id" form:"id"`
	// 服务名称
	ServiceName string `json:"service_name" form:"service_name"`
	// 服务描述
	ServiceDesc string `json:"serbice_desc" form:"serbice_desc"`
	// 类型
	LoadType int `json:"load_type" form:"load_type"`
	// 服务地址
	ServiceAddr string `json:"service_addr" form:"service_addr"`
	// QPS
	Qps int64 `json:"qps" form:"qps"`
	// 日请求总数
	Qpd int64 `json:"qpd" form:"qpd"`
	// 节点数
	TotalNode int `json:"total_node" form:"total_node"`
}

// ServiceDeleteInput ： 删除接口输入结构体
type ServiceDeleteInput struct {
	// 服务ID
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"" validate:"required"`
}

func (params *ServiceDeleteInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ServiceAddHTTPInput ： 添加列表输入结构体
type ServiceAddHTTPInput struct {
	// 服务名称
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required"`
	// 服务描述
	ServiceDesc string `json:"serbice_desc" form:"serbice_desc" comment:"服务描述" example:"" validate:"required"`
	// 接入类型
	RuleType string `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"required"`
	// 接入路径
	Rule string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required"`
	// 是否支持HTTPS, 1=支持
	NeedHTTPS int `json:"need_https" form:"need_https" comment:"是否支持HTTPS:1=支持" example:"" validate:"required"`
	// 是否启用strip_url, 1=启用
	NeedStripURL string `json:"need_strip_url" form:"need_strip_url" comment:"启用strip_url 1=启用" example:"" validate:"required"`
	// 是否支持websocket, 1=支持
	NeedSocket string `json:"need_strip_url" form:"need_strip_url" comment:"启用strip_url 1=启用" example:"" validate:"required"`
}

func (params *ServiceAddHTTPInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
