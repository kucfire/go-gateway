package dto

type PanelGroupDataOutput struct {
	ServiceNum      int64 `json:"service_num"`
	AppNum          int64 `json:"app_num"`
	CurrentQPS      int64 `json:"current_qps"`
	TodayRequestNum int64 `json:"today_request_num"`
}
