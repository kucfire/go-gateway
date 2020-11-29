package dto

type AdminInfoOutput struct {
	ID              int64 `json:"id"`
	AppNum          int64 `json:"app_num"`
	CurrentQPS      int64 `json:"current_qps"`
	TodayRequestNum int64 `json:"today_request_num"`
}
