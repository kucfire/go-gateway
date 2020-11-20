package dao

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"Info" gorm:"column:Info" description:"基本信息"`
	AccessControl *ServiceAccessControl `json:"access_control" gorm:"column:access_control" description:"基本信息"`
	HTTP          *ServiceHTTPRule      `json:"HTTP" gorm:"column:HTTP" description:"HTTP rule"`
}
