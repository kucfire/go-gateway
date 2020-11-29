package dto

type AdminInfoOutput struct {
	ID              int64 `json:"id"`
	AppNum          int64
	CurrentQPS      int64
	TodayRequestNum int64
}
