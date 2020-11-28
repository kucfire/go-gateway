package dto

import (
	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
)

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

type ServiceListOutput struct {
	// 总数
	Total string `json:"total" form:"total" comment:"总数" example:"" validate:""`
	//
	List string `json:"list" form:"list" comment:"总数" example:"" validate:""`
}

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
	ServiceAddr int `json:"service_addr" form:"service_addr"`
	// QPS
	Qps int `json:"qps" form:"qps"`
	// 日请求总数
	Qpd int `json:"qpd" form:"qpd"`
	// 节点数
	ServiceNode int `json:"service_node" form:"service_node"`
}