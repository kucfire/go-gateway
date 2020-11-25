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

// ServiceStatOutput ： 服务统计输出结构体
type ServiceStatOutput struct {
	// 昨日统计结果
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" example:"" validate:""`
	// 今日统计结果
	Today []int64 `json:"today" form:"today" comment:"今日流量" example:"" validate:""`
}

/*
	------------------------HTTP MODULE----------------------------
*/

// ServiceAddHTTPInput ： 添加列表输入结构体
type ServiceAddHTTPInput struct {
	// db.gateway_service_info
	// 服务名称
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"`
	// 服务描述
	ServiceDesc string `json:"serbice_desc" form:"serbice_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`

	// db.gateway_service_http_rule
	// 接入类型
	RuleType int `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`
	// 接入路径
	Rule string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required,vaild_rule"`
	// 是否支持HTTPS, 1=支持
	NeedHTTPS int `json:"need_https" form:"need_https" comment:"是否支持HTTPS:1=支持" example:"" validate:"max=1,min=0"`
	// 是否启用strip_uri, 1=启用
	NeedStripURI int `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri 1=启用" example:"" validate:"max=1,min=0"`
	// 是否支持websocket, 1=支持
	NeedWEBSocket int `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket 1=支持" example:"" validate:"max=1,min=0"`
	// url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
	URLRewrite string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔" example:"" validate:"valid_url_rewrite"`
	// header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔" example:"" validate:"valid_header_transfor"`

	// db.gateway_service_access_control
	// 是否开启权限 1=开启
	OpenAuth int `json:"open_auth" form:"open_auth" comment:"是否开启权限 1=开启" example:"" validate:"max=1,min=0"`
	// 黑名单ip
	BlackList string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`
	// 白名单ip
	WhiteList string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`
	// 客户端ip限流
	ClientIPFlowLimit int `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"" validate:"min=0"`
	// 服务器限流
	ServiceFlowLimit int `json:"service_flow_limit" form:"service_flow_limit" comment:"服务器限流" example:"" validate:"min=0"`

	// db.gateway_service_load_balance
	// 轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash
	RoundType int `json:"round_type" form:"round_type" comment:"轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash" example:"" validate:"max=3,min=0"`
	// ip列表
	IPList string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"required,valid_ip_list"`
	// 权重列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weight_list"`
	// 建立连接超时，单位s
	UpstreamConnectTimeout int `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s" example:"" validate:"min=0"`
	// 获取header超时，单位s
	UpstreamHeaderTimeout int `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s" example:"" validate:"min=0"`
	// 链接最大空闲时间，单位s
	UpstreamIdleTimeout int `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间，单位s" example:"" validate:"min=0"`
	// 最大空闲链接数
	UpstreamMaxIdle int `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"" validate:"min=0"`
}

func (params *ServiceAddHTTPInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ServiceAddHTTPInput ： 添加列表输入结构体
type ServiceUpdateHTTPInput struct {
	// ID
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"62" validate:"required,min=1"`

	// db.gateway_service_info
	// 服务名称
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"test_http_service_indb" validate:"required,valid_service_name"`
	// 服务描述
	ServiceDesc string `json:"serbice_desc" form:"serbice_desc" comment:"服务描述" example:"test_http_service_indb" validate:"required,max=255,min=1"`

	// db.gateway_service_http_rule
	// 接入类型
	RuleType int `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`
	// 接入路径
	Rule string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"/test_http_service_indb" validate:"required,vaild_rule"`
	// 是否支持HTTPS, 1=支持
	NeedHTTPS int `json:"need_https" form:"need_https" comment:"是否支持HTTPS:1=支持" example:"" validate:"max=1,min=0"`
	// 是否启用strip_uri, 1=启用
	NeedStripURI int `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri 1=启用" example:"" validate:"max=1,min=0"`
	// 是否支持websocket, 1=支持
	NeedWEBSocket int `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket 1=支持" example:"" validate:"max=1,min=0"`
	// url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
	URLRewrite string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔" example:"" validate:"valid_url_rewrite"`
	// header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔" example:"" validate:"valid_header_transfor"`

	// db.gateway_service_access_control
	// 是否开启权限 1=开启
	OpenAuth int `json:"open_auth" form:"open_auth" comment:"是否开启权限 1=开启" example:"" validate:"max=1,min=0"`
	// 黑名单ip
	BlackList string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`
	// 白名单ip
	WhiteList string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`
	// 客户端ip限流
	ClientIPFlowLimit int `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"" validate:"min=0"`
	// 服务器限流
	ServiceFlowLimit int `json:"service_flow_limit" form:"service_flow_limit" comment:"服务器限流" example:"" validate:"min=0"`

	// db.gateway_service_load_balance
	// 轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash
	RoundType int `json:"round_type" form:"round_type" comment:"轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash" example:"" validate:"max=3,min=0"`
	// ip列表
	IPList string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"127.0.0.1:80" validate:"required,valid_ip_list"`
	// 权重列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"50" validate:"required,valid_weight_list"`
	// 建立连接超时，单位s
	UpstreamConnectTimeout int `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s" example:"" validate:"min=0"`
	// 获取header超时，单位s
	UpstreamHeaderTimeout int `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s" example:"" validate:"min=0"`
	// 链接最大空闲时间，单位s
	UpstreamIdleTimeout int `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间，单位s" example:"" validate:"min=0"`
	// 最大空闲链接数
	UpstreamMaxIdle int `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"" validate:"min=0"`
}

func (params *ServiceUpdateHTTPInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

/*
	------------------------GRPC MODULE----------------------------
*/

// ServiceAddHTTPInput ： 添加列表输入结构体
type ServiceAddGRPCInput struct {
	// db.gateway_service_info
	// 服务名称
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"`
	// 服务描述
	ServiceDesc string `json:"serbice_desc" form:"serbice_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`

	// db.gateway_service_http_rule
	// 接入类型
	RuleType int `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`
	// 接入路径
	Rule string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required,vaild_rule"`
	// 是否支持HTTPS, 1=支持
	NeedHTTPS int `json:"need_https" form:"need_https" comment:"是否支持HTTPS:1=支持" example:"" validate:"max=1,min=0"`
	// 是否启用strip_uri, 1=启用
	NeedStripURI int `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri 1=启用" example:"" validate:"max=1,min=0"`
	// 是否支持websocket, 1=支持
	NeedWEBSocket int `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket 1=支持" example:"" validate:"max=1,min=0"`
	// url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
	URLRewrite string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔" example:"" validate:"valid_url_rewrite"`
	// header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔" example:"" validate:"valid_header_transfor"`

	// db.gateway_service_access_control
	// 是否开启权限 1=开启
	OpenAuth int `json:"open_auth" form:"open_auth" comment:"是否开启权限 1=开启" example:"" validate:"max=1,min=0"`
	// 黑名单ip
	BlackList string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`
	// 白名单ip
	WhiteList string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`
	// 客户端ip限流
	ClientIPFlowLimit int `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"" validate:"min=0"`
	// 服务器限流
	ServiceFlowLimit int `json:"service_flow_limit" form:"service_flow_limit" comment:"服务器限流" example:"" validate:"min=0"`

	// db.gateway_service_load_balance
	// 轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash
	RoundType int `json:"round_type" form:"round_type" comment:"轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash" example:"" validate:"max=3,min=0"`
	// ip列表
	IPList string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"required,valid_ip_list"`
	// 权重列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weight_list"`
	// 建立连接超时，单位s
	UpstreamConnectTimeout int `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s" example:"" validate:"min=0"`
	// 获取header超时，单位s
	UpstreamHeaderTimeout int `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s" example:"" validate:"min=0"`
	// 链接最大空闲时间，单位s
	UpstreamIdleTimeout int `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间，单位s" example:"" validate:"min=0"`
	// 最大空闲链接数
	UpstreamMaxIdle int `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"" validate:"min=0"`
}

func (params *ServiceAddGRPCInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ServiceAddHTTPInput ： 添加列表输入结构体
type ServiceUpdateGRPCInput struct {
	// ID
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"62" validate:"required,min=1"`

	// db.gateway_service_info
	// 服务名称
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"test_http_service_indb" validate:"required,valid_service_name"`
	// 服务描述
	ServiceDesc string `json:"serbice_desc" form:"serbice_desc" comment:"服务描述" example:"test_http_service_indb" validate:"required,max=255,min=1"`

	// db.gateway_service_http_rule
	// 接入类型
	RuleType int `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`
	// 接入路径
	Rule string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"/test_http_service_indb" validate:"required,vaild_rule"`
	// 是否支持HTTPS, 1=支持
	NeedHTTPS int `json:"need_https" form:"need_https" comment:"是否支持HTTPS:1=支持" example:"" validate:"max=1,min=0"`
	// 是否启用strip_uri, 1=启用
	NeedStripURI int `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri 1=启用" example:"" validate:"max=1,min=0"`
	// 是否支持websocket, 1=支持
	NeedWEBSocket int `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket 1=支持" example:"" validate:"max=1,min=0"`
	// url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
	URLRewrite string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔" example:"" validate:"valid_url_rewrite"`
	// header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔" example:"" validate:"valid_header_transfor"`

	// db.gateway_service_access_control
	// 是否开启权限 1=开启
	OpenAuth int `json:"open_auth" form:"open_auth" comment:"是否开启权限 1=开启" example:"" validate:"max=1,min=0"`
	// 黑名单ip
	BlackList string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`
	// 白名单ip
	WhiteList string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`
	// 客户端ip限流
	ClientIPFlowLimit int `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"" validate:"min=0"`
	// 服务器限流
	ServiceFlowLimit int `json:"service_flow_limit" form:"service_flow_limit" comment:"服务器限流" example:"" validate:"min=0"`

	// db.gateway_service_load_balance
	// 轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash
	RoundType int `json:"round_type" form:"round_type" comment:"轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash" example:"" validate:"max=3,min=0"`
	// ip列表
	IPList string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"127.0.0.1:80" validate:"required,valid_ip_list"`
	// 权重列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"50" validate:"required,valid_weight_list"`
	// 建立连接超时，单位s
	UpstreamConnectTimeout int `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s" example:"" validate:"min=0"`
	// 获取header超时，单位s
	UpstreamHeaderTimeout int `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s" example:"" validate:"min=0"`
	// 链接最大空闲时间，单位s
	UpstreamIdleTimeout int `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间，单位s" example:"" validate:"min=0"`
	// 最大空闲链接数
	UpstreamMaxIdle int `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"" validate:"min=0"`
}

func (params *ServiceUpdateHTTPInput) ServiceUpdateGRPCInput(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

/*
	------------------------TCP MODULE----------------------------
*/

// ServiceAddHTTPInput ： 添加列表输入结构体
type ServiceAddTCPInput struct {
	// db.gateway_service_info
	// 服务名称
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"`
	// 服务描述
	ServiceDesc string `json:"serbice_desc" form:"serbice_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`

	// db.gateway_service_http_rule
	// 接入类型
	RuleType int `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`
	// 接入路径
	Rule string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required,vaild_rule"`
	// 是否支持HTTPS, 1=支持
	NeedHTTPS int `json:"need_https" form:"need_https" comment:"是否支持HTTPS:1=支持" example:"" validate:"max=1,min=0"`
	// 是否启用strip_uri, 1=启用
	NeedStripURI int `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri 1=启用" example:"" validate:"max=1,min=0"`
	// 是否支持websocket, 1=支持
	NeedWEBSocket int `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket 1=支持" example:"" validate:"max=1,min=0"`
	// url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
	URLRewrite string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔" example:"" validate:"valid_url_rewrite"`
	// header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔" example:"" validate:"valid_header_transfor"`

	// db.gateway_service_access_control
	// 是否开启权限 1=开启
	OpenAuth int `json:"open_auth" form:"open_auth" comment:"是否开启权限 1=开启" example:"" validate:"max=1,min=0"`
	// 黑名单ip
	BlackList string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`
	// 白名单ip
	WhiteList string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`
	// 客户端ip限流
	ClientIPFlowLimit int `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"" validate:"min=0"`
	// 服务器限流
	ServiceFlowLimit int `json:"service_flow_limit" form:"service_flow_limit" comment:"服务器限流" example:"" validate:"min=0"`

	// db.gateway_service_load_balance
	// 轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash
	RoundType int `json:"round_type" form:"round_type" comment:"轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash" example:"" validate:"max=3,min=0"`
	// ip列表
	IPList string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"required,valid_ip_list"`
	// 权重列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weight_list"`
	// 建立连接超时，单位s
	UpstreamConnectTimeout int `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s" example:"" validate:"min=0"`
	// 获取header超时，单位s
	UpstreamHeaderTimeout int `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s" example:"" validate:"min=0"`
	// 链接最大空闲时间，单位s
	UpstreamIdleTimeout int `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间，单位s" example:"" validate:"min=0"`
	// 最大空闲链接数
	UpstreamMaxIdle int `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"" validate:"min=0"`
}

func (params *ServiceAddTCPInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

// ServiceAddHTTPInput ： 添加列表输入结构体
type ServiceUpdateTCPInput struct {
	// ID
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"62" validate:"required,min=1"`

	// db.gateway_service_info
	// 服务名称
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"test_http_service_indb" validate:"required,valid_service_name"`
	// 服务描述
	ServiceDesc string `json:"serbice_desc" form:"serbice_desc" comment:"服务描述" example:"test_http_service_indb" validate:"required,max=255,min=1"`

	// db.gateway_service_http_rule
	// 接入类型
	RuleType int `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`
	// 接入路径
	Rule string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"/test_http_service_indb" validate:"required,vaild_rule"`
	// 是否支持HTTPS, 1=支持
	NeedHTTPS int `json:"need_https" form:"need_https" comment:"是否支持HTTPS:1=支持" example:"" validate:"max=1,min=0"`
	// 是否启用strip_uri, 1=启用
	NeedStripURI int `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri 1=启用" example:"" validate:"max=1,min=0"`
	// 是否支持websocket, 1=支持
	NeedWEBSocket int `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket 1=支持" example:"" validate:"max=1,min=0"`
	// url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔
	URLRewrite string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔" example:"" validate:"valid_url_rewrite"`
	// header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换支持增加(add)、删除(del)、修改(edit) 格式：add headname headvalue 多个逗号间隔" example:"" validate:"valid_header_transfor"`

	// db.gateway_service_access_control
	// 是否开启权限 1=开启
	OpenAuth int `json:"open_auth" form:"open_auth" comment:"是否开启权限 1=开启" example:"" validate:"max=1,min=0"`
	// 黑名单ip
	BlackList string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`
	// 白名单ip
	WhiteList string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`
	// 客户端ip限流
	ClientIPFlowLimit int `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流" example:"" validate:"min=0"`
	// 服务器限流
	ServiceFlowLimit int `json:"service_flow_limit" form:"service_flow_limit" comment:"服务器限流" example:"" validate:"min=0"`

	// db.gateway_service_load_balance
	// 轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash
	RoundType int `json:"round_type" form:"round_type" comment:"轮询方式 0=random 1=round_robin 2=wieght_round_robin 3=ip_hash" example:"" validate:"max=3,min=0"`
	// ip列表
	IPList string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"127.0.0.1:80" validate:"required,valid_ip_list"`
	// 权重列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"50" validate:"required,valid_weight_list"`
	// 建立连接超时，单位s
	UpstreamConnectTimeout int `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s" example:"" validate:"min=0"`
	// 获取header超时，单位s
	UpstreamHeaderTimeout int `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s" example:"" validate:"min=0"`
	// 链接最大空闲时间，单位s
	UpstreamIdleTimeout int `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间，单位s" example:"" validate:"min=0"`
	// 最大空闲链接数
	UpstreamMaxIdle int `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"" validate:"min=0"`
}

func (params *ServiceUpdateTCPInput) ServiceUpdateGRPCInput(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
