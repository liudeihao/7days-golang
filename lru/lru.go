package lru

import "container/list"

type Cache struct {
    maxBytes  int64 // =0表示不设限
    nbytes    int64
    ll        *list.List
    cache     map[string]*list.Element
    OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
    return &Cache{
        maxBytes:  maxBytes,
        ll:        list.New(),
        cache:     make(map[string]*list.Element),
        OnEvicted: onEvicted,
    }
}

func (c *Cache) Get(key string) (value Value, ok bool) {
    e, ok := c.cache[key]
    if !ok {
        return nil, false
    }
    c.ll.MoveToBack(e)
    kv := e.Value.(*entry)
    return kv.value, true
}

func (c *Cache) Add(key string, value Value) {
    if e, ok := c.cache[key]; ok {
        c.ll.MoveToBack(e)
        e.Value.(*entry).value = value
        return
    } else {
        c.cache[key] = c.ll.PushBack(&entry{key: key, value: value})
        c.nbytes += int64(len(key)) + int64(value.Len())
    }
    for c.maxBytes != 0 && c.nbytes > c.maxBytes {
        c.RemoveOldest()
    }
}

func (c *Cache) RemoveOldest() {
    e := c.ll.Front()
    if e != nil {
        v := c.ll.Remove(e)
        kv := v.(*entry)
        delete(c.cache, kv.key)
        c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
        if c.OnEvicted != nil {
            c.OnEvicted(kv.key, kv.value)
        }
    }
}

type entry struct {
    key   string
    value Value
}

type Value interface {
    Len() int // 返回占用内存大小
}

func (c *Cache) Len() int {
    // 添加了多少条数据
    return c.ll.Len()
}
