package dao

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"Info" gorm:"column:Info" description:"基本信息"`
	AccessControl *ServiceAccessControl `json:"access_control" gorm:"column:access_control" description:"access control"`
	HTTP          *ServiceHTTPRule      `json:"http" gorm:"column:http" description:"http rule"`
	GRPC          *ServiceGRPCRule      `json:"gprc" gorm:"column:grpc" description:"gprc rule"`
}
