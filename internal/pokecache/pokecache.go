package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu    sync.RWMutex
	cache map[string]cacheEntry
}

func (p *Cache) Add(key string, value []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (p *Cache) Get(key string) ([]byte, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.cache[key]
	if !ok {
		return nil, false
	}
	return val.val, true
}

func (p *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		for key, entry := range p.cache {
			if time.Since(entry.createdAt) > interval {
				delete(p.cache, key)
			}
		}
		p.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cache: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}
