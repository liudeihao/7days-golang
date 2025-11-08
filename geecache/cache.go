package geecache

import (
    "geecache/geecache/lru"
    "sync"
)

type cache struct {
    lru        *lru.Cache
    mu         sync.Mutex
    cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.lru == nil {
        // 延迟初始化
        c.lru = lru.New(c.cacheBytes, nil)
    }
    c.lru.Add(key, value)
    c.cacheBytes += int64(len(key))
}

func (c *cache) get(key string) (ByteView, bool) {
    c.mu.Lock()
    var value ByteView
    defer c.mu.Unlock()
    if c.lru == nil {
        return value, false
    }
    if v, ok := c.lru.Get(key); ok {
        return v.(ByteView), ok
    }
    return value, false
}
