package consistentHashMapBalance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Hash func(data []byte) uint32

type UInt32Slice []uint32

// get the length of UInt32Slice
func (s UInt32Slice) Len() int {
	return len(s)
}

// get the judge about two element
func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashMapBalance struct {
	mux     sync.RWMutex
	hash    Hash
	replace int               // 复制因子
	keys    UInt32Slice       // 已排序的节点hash切片
	hashMap map[uint32]string // 节点哈希和key的map，键是hash值，值是节点key
}

func NewConsistentHashBanlance(replicas int, fn Hash) *ConsistentHashMapBalance {
	m := &ConsistentHashMapBalance{
		replace: replicas,
		hash:    fn,
		hashMap: make(map[uint32]string),
	}

	// if hash is nil, we will set a default value
	if m.hash == nil {
		// 最多32位，保证是一个2^32-1
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

func (c *ConsistentHashMapBalance) Add(params ...string) error {
	// 零值处理
	if len(params) == 0 {
		return errors.New("param need a element at least")
	}

	// 过滤其它输入的值，只取第一个值
	addr := params[0]

	c.mux.Lock()
	defer c.mux.Unlock()

	// 结合复制因子计算所有虚拟节点的hash值，并存入M.keys中，同时在m.hashMap中保存哈希值和key的映射
	for i := 0; i < c.replace; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr)) // 计算哈希值
		c.keys = append(c.keys, hash)                  // 存入哈希值
		c.hashMap[hash] = addr                         // key-value 进行值映射
	}

	sort.Sort(c.keys) //进行排序，方便进行二分查找获取
	return nil
}

func (c *ConsistentHashMapBalance) Get(key string) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("node is empty")
	}

	hash := c.hash([]byte(key))

	// 通过二分查找获取最后节点, 第一个c.keys的值大于key的hash值就是目标节点
	idx := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })

	if idx == len(c.keys) {
		idx = 0
	}
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.hashMap[c.keys[idx]], nil
}

func (c *ConsistentHashMapBalance) IsEmpty() bool {
	return len(c.keys) == 0
}
