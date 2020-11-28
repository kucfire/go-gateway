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

// ServiceListOutput ： 服务列表输出结构体
type AppListOutput struct {
	// 总数
	Total int64 `json:"total" form:"total" comment:"总数" example:"" validate:""`
	// 列表
	List []AppListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`
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
	Sercet string `json:"sercet" form:"sercet"`
	// 服务地址
	WhiteIPS string `json:"white_ips" form:"white_ips"`
	// QPS
	QPS int64 `json:"qps" form:"qps"`
	// 日请求总数
	QPD int64 `json:"qpd" form:"qpd"`
}
