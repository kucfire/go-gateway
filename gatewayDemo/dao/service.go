package dao

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"info" gorm:"column:Info" description:"基本信息"`
	AccessControl *ServiceAccessControl `json:"access_control" gorm:"column:access_control" description:"access control"`
	HTTPRule      *ServiceHTTPRule      `json:"http" gorm:"column:http" description:"http rule"`
	GRPCRule      *ServiceGRPCRule      `json:"gprc" gorm:"column:grpc" description:"gprc rule"`
	LoadBalance   *ServiceLoadBalance   `json:"load_balance" gorm:"column:load_balance" description:"load balance"`
	TCPRule       *ServiceTCPRule       `json:"tcp" gorm:"column:tcp" description:"tcp rule"`
}
