package loadbalance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// Author:Boyn
// Date:2020/8/31

type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashBalance struct {
	mux      sync.RWMutex
	hash     Hash
	replicas int               //复制因子
	keys     UInt32Slice       //已排序的节点hash切片
	hashMap  map[uint32]string //节点哈希和Key的map,键是hash值，值是节点key
}

const replicas = 3

func NewConsistentHash() *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     crc32.ChecksumIEEE,
		hashMap:  make(map[uint32]string),
	}
	return m
}

func (c *ConsistentHashBalance) Add(params ...ConfigValue) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0].Value
	c.mux.Lock()
	defer c.mux.Unlock()
	// 结合复制因子计算所有虚拟节点的hash值，并存入m.keys中，同时在m.hashMap中保存哈希值和key的映射
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr
	}
	// 对所有虚拟节点的哈希值进行排序，方便之后进行二分查找
	sort.Sort(c.keys)
	return nil
}
func (c *ConsistentHashBalance) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashBalance) Get(key string) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("node is empty")
	}
	hash := c.hash([]byte(key))

	// 通过二分查找获取最优节点，第一个"服务器hash"值大于"数据hash"值的就是最优"服务器节点"
	idx := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })

	// 如果查找结果 大于 服务器节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，那么放入第一个节点中
	if idx == len(c.keys) {
		idx = 0
	}
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.hashMap[c.keys[idx]], nil
}

func (c *ConsistentHashBalance) GetAll() []string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	addrs := make([]string, 0)
	for _, v := range c.hashMap {
		addrs = append(addrs, v)
	}
	return addrs
}

func (c *ConsistentHashBalance) Update(newRss []ConfigValue) error {
	newHash := NewConsistentHash()
	for _, v := range newRss {
		err := newHash.Add(v)
		if err != nil {
			return err
		}
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	c.keys = newHash.keys
	c.hashMap = newHash.hashMap
	return nil
}
