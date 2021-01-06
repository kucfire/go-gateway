package public

const (
	ValidatorKey        = "ValidatorKey"
	TranslatorKey       = "TranslatorKey"
	AdminSessionInfoKey = "AdminSessionInfo"

	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotalPrefix   = "flow_total"
	FlowServicePrefix = "flow_Service_"
	FlowAppPrefix     = "flow_app_"

	JwtSignKey = "my_sign_key"
	JwtExpires = 60 * 60
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeGRPC: "GRPC",
		LoadTypeTCP:  "TCP",
	}
)
