package dao

import (
	"errors"
	"gatewayDemo/dto"
	"gatewayDemo/public"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

var ServiceManagerHandler *ServiceManager

func init() {
	ServiceManagerHandler = NewServiceManager()
}

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"info" gorm:"column:Info" description:"基本信息"`
	AccessControl *ServiceAccessControl `json:"access_control" gorm:"column:access_control" description:"access control"`
	HTTPRule      *ServiceHTTPRule      `json:"http" gorm:"column:http" description:"http rule"`
	GRPCRule      *ServiceGRPCRule      `json:"grpc" gorm:"column:grpc" description:"gprc rule"`
	LoadBalance   *ServiceLoadBalance   `json:"load_balance" gorm:"column:load_balance" description:"load balance"`
	TCPRule       *ServiceTCPRule       `json:"tcp" gorm:"column:tcp" description:"tcp rule"`
}

type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker       sync.RWMutex
	init         sync.Once
	errMsg       error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
}

func (s *ServiceManager) GetTCPServiceList() []*ServiceDetail {
	list := []*ServiceDetail{}
	for _, servicelist := range s.ServiceSlice {
		temp := servicelist
		if temp.Info.LoadType == public.LoadTypeTCP {
			list = append(list, temp)
		}
	}
	return list
}

func (s *ServiceManager) GetGRPCServiceList() []*ServiceDetail {
	list := []*ServiceDetail{}
	for _, servicelist := range s.ServiceSlice {
		temp := servicelist
		if temp.Info.LoadType == public.LoadTypeGRPC {
			list = append(list, temp)
		}
	}
	return list
}

func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*ServiceDetail, error) {

	// 前缀匹配 ： /abc ==> serviceSlice.rule
	// 域名匹配 ： www.test.com == > serviceSLice.Rule

	// host c.Request.Host
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	// fmt.Println("host : ", host)

	// path c.Request.URL.Path
	path := c.Request.URL.Path
	// fmt.Println("path : ", path)

	for _, serviceItem := range s.ServiceSlice {
		// 非HTTP的类型则直接跳过
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			continue
		}

		// 判断匹配类型是否为域名
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}

		// 判断匹配类型是否为url前缀
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			// fmt.Println("rule : ", serviceItem.HTTPRule.Rule)
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}

	return nil, errors.New("not matched service")
}

func (s *ServiceManager) LoadOnce() error {
	s.init.Do(func() {
		serviceInfo := &ServiceInfo{}

		// 设置*gin.context
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		// 连接池
		tx, err := lib.GetGormPool("default")
		if err != nil {
			s.errMsg = err
			return
		}

		// 从DB中分页读取基本信息
		// 取出所有数据
		params := &dto.ServiceListInput{
			PageSize: 99999,
			PageNo:   1,
		}
		list, _, err := serviceInfo.PageList(c, tx, params)
		if err != nil {
			s.errMsg = err
			return
		}

		// 遍历整个结果列表
		s.Locker.Lock()
		defer s.Locker.Unlock()
		for _, listItem := range list {
			tmp := listItem
			serviceDetail, err := tmp.ServiceDetail(c, tx, &tmp)
			if err != nil {
				s.errMsg = err
				return
			}
			// fmt.Println("serviceDetail")
			// fmt.Println(public.ObjToJson(serviceDetail))
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
	})
	return s.errMsg
}
