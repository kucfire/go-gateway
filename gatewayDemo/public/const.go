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
	FlowServicePrefix = "flow_Service"
	FlowAppPrefix     = "flow_app"
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeGRPC: "GRPC",
		LoadTypeTCP:  "TCP",
	}
)