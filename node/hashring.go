package node

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type ConsistentHashRing struct {
	sync.RWMutex
	nodes        map[int]string
	sortedHashes []int
	replicas     int
}

func NewConsistentHashRing(replicas int) *ConsistentHashRing {
	return &ConsistentHashRing{
		nodes:        make(map[int]string),
		sortedHashes: []int{},
		replicas:     replicas,
	}
}

func (c *ConsistentHashRing) AddNode(node string) {
	c.Lock()
	defer c.Unlock()

	for i := 0; i < c.replicas; i++ {
		hash := int(hashKey(node + strconv.Itoa(i)))
		c.nodes[hash] = node
		c.sortedHashes = append(c.sortedHashes, hash)
	}
	sort.Ints(c.sortedHashes)
}

func (c *ConsistentHashRing) GetNode(key string) string {
	c.RLock()
	defer c.RUnlock()

	if len(c.nodes) == 0 {
		return ""
	}

	hash := int(hashKey(key))
	idx := sort.Search(len(c.sortedHashes), func(i int) bool {
		return c.sortedHashes[i] >= hash
	})

	if idx == len(c.sortedHashes) {
		idx = 0
	}
	return c.nodes[c.sortedHashes[idx]]
}

func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *ConsistentHashRing) AddNodes(peers []string) {
	for _, peer := range peers {
		c.AddNode(peer)
	}
}
