package core

import (
	"liteGoCache/policy"
	"sync"
)

type goCache struct {
	mu         sync.Mutex
	lru        *policy.LRUCache
	cacheBytes int64
}

func (c *goCache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = policy.NewLRUCache(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *goCache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
