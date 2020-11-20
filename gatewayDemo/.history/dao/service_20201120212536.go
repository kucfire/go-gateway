package dao

type ServiceDetail struct {
	Info int64 `json:"Info" gorm:"column:Info" description:"服务类型"`
}
