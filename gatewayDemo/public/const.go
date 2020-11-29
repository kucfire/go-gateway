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
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeGRPC: "GRPC",
		LoadTypeTCP:  "TCP",
	}
)
