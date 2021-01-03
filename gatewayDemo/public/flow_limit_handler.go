package public

import (
	"sync"

	"golang.org/x/time/rate"
)

var FlowLimiterHandler *FlowLimiter

func init() {
	FlowLimiterHandler = NewFlowLimiterService()
}

type FlowLimiter struct {
	// 服务较多时，可以用
	FlowLimiterMap   map[string]*FlowLimiterItem
	FlowLimiterSlice []*FlowLimiterItem
	Locker           sync.RWMutex
	// init             sync.Once
	// errMsg           error
}

type FlowLimiterItem struct {
	ServiceName string
	Limiter     *rate.Limiter
}

func NewFlowLimiterService() *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap:   map[string]*FlowLimiterItem{},
		FlowLimiterSlice: []*FlowLimiterItem{},
		Locker:           sync.RWMutex{},
	}
}

func (limit *FlowLimiter) GetFlowLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	for _, item := range limit.FlowLimiterSlice {
		if item.ServiceName == serviceName {
			return item.Limiter, nil
		}
	}

	// 超过三倍的时候直接熔断
	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	newItem := &FlowLimiterItem{
		ServiceName: serviceName,
		Limiter:     newLimiter,
	}

	limit.FlowLimiterSlice = append(limit.FlowLimiterSlice, newItem)
	limit.Locker.Lock()
	defer limit.Locker.Unlock()
	limit.FlowLimiterMap[serviceName] = newItem

	return newLimiter, nil
}
