package policy

import (
	"container/list"
)

type LRUCache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func NewLRUCache(maxBytes int64, onEvicted func(string, Value)) *LRUCache {
	return &LRUCache{
		maxBytes:  maxBytes,
		nbytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get Func
func (lc *LRUCache) Get(key string) (value Value, ok bool) {
	if ele, ok := lc.cache[key]; ok {
		lc.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// Delete Func
func (lc *LRUCache) RemoveOldest() {
	ele := lc.ll.Back()
	if ele != nil {
		lc.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(lc.cache, kv.key)
		lc.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if lc.OnEvicted != nil {
			lc.OnEvicted(kv.key, kv.value)
		}
	}
}

// append/update
func (lc *LRUCache) Add(key string, value Value) {
	if ele, ok := lc.cache[key]; ok {
		lc.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		lc.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := lc.ll.PushFront(&entry{key, value})
		lc.cache[key] = ele
		lc.nbytes += int64(len(key)) + int64(value.Len())
	}
	for lc.maxBytes != 0 && lc.maxBytes < lc.nbytes {
		lc.RemoveOldest()
	}
}

// back element num
func (lc *LRUCache) Len() int {
	return lc.ll.Len()
}
