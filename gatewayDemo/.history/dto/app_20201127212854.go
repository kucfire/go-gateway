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

// ServiceListOutput ： 租户列表输出结构体
type AppListOutput struct {
	// 总数
	Total int64 `json:"total" form:"total" comment:"总数" example:"" validate:""`
	// 列表
	List []AppListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`
}

// AppListItemOutput ： 租户列表中的表单输出结构体
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

// AppAddInput : 添加租户app
type AppAddInput struct {
	// 服务名称
	AppID string `json:"app_id" form:"app_id" comment:"租户id" example:"" validate:"required"`
	// 服务名称
	Name string `json:"name" form:"name" comment:"租户名称" example:"" validate:"required"`
	// 服务名称
	Secret string `json:"secret" form:"secret" comment:"密钥" example:"" validate:"required"`
	// 服务名称
	WhiteIPS string `json:"white_ips" form:"white_ips" comment:"ip白名单，支持前缀匹配" example:"" validate:"required"`
}

func (params *AppAddInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// AppUpdateInput : 添加租户app
type AppUpdateInput struct {
	// 服务名称:需要ID进行校验
	AppID string `json:"app_id" form:"app_id" comment:"租户id" example:"" validate:"required"`
	// 服务名称
	Name string `json:"name" form:"name" comment:"租户名称" example:"" validate:"required"`
	// 服务名称
	Secret string `json:"secret" form:"secret" comment:"密钥" example:"" validate:"required,vaild_secret"`
	// 服务名称
	WhiteIPS string `json:"white_ips" form:"white_ips" comment:"ip白名单，支持前缀匹配" example:"" validate:"required,vaild_white_ips"`
}

func (params *AppUpdateInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// AppDeleteInput : 删除租户app
type AppDeleteInput struct {
	// 服务名称:需要ID进行校验
	AppID string `json:"app_id" form:"app_id" comment:"租户id" example:"" validate:"required"`
}

func (params *AppDeleteInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// AppStatInput : 租户流量统计
type AppStatInput struct {
	// 昨日统计结果
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" example:"" validate:""`
	// 今日统计结果
	Today []int64 `json:"today" form:"today" comment:"今日流量" example:"" validate:""`
}

func (params *AppStatInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
