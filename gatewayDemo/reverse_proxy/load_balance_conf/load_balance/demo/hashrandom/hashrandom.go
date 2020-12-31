package hashrandom

import (
	"errors"
	"fmt"
	"gatewayDemo/reverse_proxy/load_balance_conf/config"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Hash func(data []byte) uint32

type Uint32Slice []uint32

func (s Uint32Slice) Len() int {
	return len(s)
}

func (s Uint32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Uint32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashmapBalance struct {
	mux     sync.RWMutex
	hash    Hash
	replace int               //复制因子
	keys    Uint32Slice       // 已排序的节点hash切片
	hashMap map[uint32]string // 节点哈希和key的map，key是hash值，value是addr

	// 观察主体
	conf config.LoadBalanceConf
}

// NewConsistentHashmapBalance : 构建一个新的对象
func NewConsistentHashmapBalance(replicas int, fn Hash) *ConsistentHashmapBalance {
	m := &ConsistentHashmapBalance{
		replace: replicas,
		hash:    fn,
		hashMap: make(map[uint32]string),
	}

	// fn为空时，hash设定为默认值
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

func (r *ConsistentHashmapBalance) Add(key ...string) error {
	if len(key) == 0 {
		return errors.New("key need 1 at least")
	}

	addr := key[0]

	r.mux.Lock()
	defer r.mux.Unlock()

	for i := 0; i < r.replace; i++ {
		hash := r.hash([]byte(strconv.Itoa(i) + addr))
		r.keys = append(r.keys, hash)
		r.hashMap[hash] = addr
	}

	sort.Sort(r.keys)
	return nil
}

func (r *ConsistentHashmapBalance) Get(key string) (string, error) {
	// 控制判断
	if r.IsEmpty() {
		return "", errors.New("hash node is empty")
	}

	// 获取key的hash值
	hash := r.hash([]byte(key))

	// 二分查找距离hash最接近的值的下标
	idx := sort.Search(len(r.keys),
		func(i int) bool {
			return r.keys[i] >= hash
		})

	// 末值判断
	if idx == len(r.keys) {
		idx = 0
	}

	// 读锁
	r.mux.RLock()
	defer r.mux.RUnlock()

	return r.hashMap[r.keys[idx]], nil
}

func (c *ConsistentHashmapBalance) SetConf(conf config.LoadBalanceConf) {
	c.conf = conf
}

func (r *ConsistentHashmapBalance) IsEmpty() bool {
	return len(r.keys) == 0
}

func (r *ConsistentHashmapBalance) Next() string {
	return ""
}

func (r *ConsistentHashmapBalance) Update() {
	// 已注册在zk集群上
	// if conf, ok := r.conf.(*config.LoadBalanceZkConf); ok {
	// 	fmt.Println("ConsistentHashmapBalance get conf : ", conf.GetConf())
	// 	r.keys = nil
	// 	r.hashMap = nil
	// 	for _, ip := range conf.GetConf() {
	// 		r.Add(strings.Split(ip, ",")...)
	// 	}
	// }

	if conf, ok := r.conf.(*config.LoadBalanceZkCheckConf); ok {
		fmt.Println("ConsistentHashmapBalance get conf : ", conf.GetConf())
		r.keys = nil
		r.hashMap = map[uint32]string{}
		for _, ip := range conf.GetConf() {
			r.Add(strings.Split(ip, ",")...)
		}
	}
}
