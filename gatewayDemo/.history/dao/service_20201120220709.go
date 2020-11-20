package dao

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"Info" gorm:"column:Info" description:"基本信息"`
	AccessControl *ServiceAccessControl `json:"HTTP" gorm:"column:HTTP" description:"基本信息"`
	HTTP          *ServiceHTTPRule      `json:"HTTP" gorm:"column:HTTP" description:"HTTP rule"`
}
