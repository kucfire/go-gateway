package weightroundrobin

import (
	"errors"
	"strconv"
)

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int
}

type WeightNode struct {
	addr            string
	weight          int // 节点权重值
	curweight       int // 临时权重值
	effectiveweight int // 有效权重值
}

// need two element ,first is IP, second is weight
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
	node.effectiveweight = node.weight
	r.rss = append(r.rss, node)
	return nil
}

func (r *WeightRoundRobinBalance) Get(key string) (string, error) {
	return "empty function", nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		// 统计节点的权重之和,在该函数结尾时需要该参数
		total += w.effectiveweight

		// 更改临时权重的值为临时权重+有效权重
		w.curweight += w.effectiveweight

		// 有效权重默认与weight相同，通讯异常时-1，通讯成功是+1，直到回复与weight大小相同为止
		if w.effectiveweight < w.weight {
			w.effectiveweight++
		}

		// 选择最大临时权重节点
		if best == nil || w.curweight > best.curweight {
			best = w
		}

	}
	if best == nil {
		return ""
	}

	best.curweight -= total
	return best.addr
}
