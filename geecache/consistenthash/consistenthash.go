package consistenthash

import (
    "hash/crc32"
    "slices"
    "strconv"
)

type Hash func(data []byte) uint32

type Map struct {
    hash     Hash
    replicas int            // 虚拟节点的倍数
    keys     []int          // 哈希环
    hashMap  map[int]string // 虚拟节点的哈希值 -> 物理节点的名称
}

func New(replicas int, fn Hash) *Map {
    m := &Map{
        replicas: replicas,
        hash:     fn,
        hashMap:  make(map[int]string),
    }
    if m.hash == nil {
        m.hash = crc32.ChecksumIEEE
    }
    return m
}

// Add 添加若干个真实节点
func (m *Map) Add(keys ...string) {
    for _, key := range keys {
        for i := 0; i < m.replicas; i++ {
            virtName := strconv.Itoa(i) + key // 虚拟节点的名称
            hash := int(m.hash([]byte(virtName)))
            m.keys = append(m.keys, hash)
            m.hashMap[hash] = key
        }
    }
    slices.Sort(m.keys)
}

func (m *Map) Get(key string) string {
    if len(m.keys) == 0 {
        return ""
    }
    hash := int(m.hash([]byte(key)))
    idx, _ := slices.BinarySearch(m.keys, hash) // >= hash
    return m.hashMap[m.keys[idx%len(m.keys)]]
}
