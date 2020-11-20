package dao

type ServiceDetail struct {
	Info        int64  `json:"load_type" gorm:"column:load_type" description:"服务类型"`
	LoadType    int    `json:"load_type" gorm:"column:load_type" description:"服务类型"`
	ServiceName string `json:"service_name" gorm:"column:service_name" description:"服务名称"`
}
