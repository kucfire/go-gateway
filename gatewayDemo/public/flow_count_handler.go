package public

import (
	"sync"
	"time"
)

var FlowCounterHandler *FlowCounter

func init() {
	FlowCounterHandler = NewFlowCountService()
}

type FlowCounter struct {
	// 服务较多时，可以用
	RedisFlowCounterMap   map[string]*RedisFlowCountService
	RedisFlowCounterSlice []*RedisFlowCountService
	Locker                sync.RWMutex
	// init             sync.Once
	// errMsg           error
}

func NewFlowCountService() *FlowCounter {
	return &FlowCounter{
		RedisFlowCounterMap:   map[string]*RedisFlowCountService{},
		RedisFlowCounterSlice: []*RedisFlowCountService{},
		Locker:                sync.RWMutex{},
	}
}

func (counter *FlowCounter) GetFlowCounter(serviceName string) (*RedisFlowCountService, error) {
	for _, item := range counter.RedisFlowCounterSlice {
		if item.AppID == serviceName {
			return item, nil
		}
	}

	// 获取redis临时存储的数据
	newCounter := NewRedisFlowCountService(serviceName, time.Second)

	counter.RedisFlowCounterSlice = append(counter.RedisFlowCounterSlice, newCounter)
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.RedisFlowCounterMap[serviceName] = newCounter

	return newCounter, nil
}
