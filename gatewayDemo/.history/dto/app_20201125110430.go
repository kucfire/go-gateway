package dto

import (
	"gatewayDemo/public"

	"github.com/gin-gonic/gin"
)

// AppListInput : 租户列表输入结构体
type AppListInput struct {
	// 关键词
	Info string `json:"info" form:"info" comment:"关键词" example:"" validate:""`
	// 页数
	PageNo int `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`
	// 每页条数
	PageSize int `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"`
}

func (params *AppListInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// AppListItemOutput ： 服务列表中的表单输出结构体
type AppListItemOutput struct {
	// 自增ID
	ID int64 `json:"id" form:"id"`
	// 租户id
	AppID string `json:"app_id" form:"app_id"`
	// 服务描述
	Name string `json:"name" form:"name"`
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