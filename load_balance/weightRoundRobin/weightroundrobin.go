package weightRoundRobin

import (
	"errors"
	"strconv"
)

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int

	// 观察主体
}

type WeightNode struct {
	addr            string
	weight          int // 节点权重
	currentWeight   int // 当前权重值
	effectiveWeight int // 有效权重值
}

func (r *WeightRoundRobinBalance) Add(key ...string) error {
	if len(key) < 2 {
		return errors.New("key need 2 at least")
	}

	parInt, err := strconv.ParseInt(key[1], 10, 64)
	if err != nil {
		return err
	}

	node := &WeightNode{
		addr:   key[0],
		weight: int(parInt),
	}
	node.effectiveWeight = node.weight
	r.rss = append(r.rss, node)

	return nil
}

func (r *WeightRoundRobinBalance) Get(key string) (string, error) {
	return "", nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0 // sum effectiveWeight
	var best *WeightNode

	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		// 统计有效权重之和
		total += w.effectiveWeight

		// 变更节点临时权重为节点的临时权重+有效权重
		w.currentWeight += w.effectiveWeight

		// 有效权重默认与权重相同，通讯异常时-1，通讯成功+1， 直到回复到weight大小
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}

		// 选择最大临时权重节点
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}

	if best == nil {
		return ""
	}

	best.currentWeight -= total
	return best.addr
}
