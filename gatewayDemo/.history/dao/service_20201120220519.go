package dao

type ServiceDetail struct {
	Info ServiceInfo `json:"Info" gorm:"column:Info" description:"基本信息"`
}
