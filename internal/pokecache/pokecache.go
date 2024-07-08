package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries              map[string]CacheEntry
	Mu                   *sync.Mutex
	InvalidationInterval time.Duration
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c Cache) ReapLoop() {
	ticker := time.NewTicker(c.InvalidationInterval)
	for {
		<-ticker.C
		c.Mu.Lock()
		for key, val := range c.Entries {
			if time.Now().After(val.createdAt) {
				delete(c.Entries, key)
			}
		}
		c.Mu.Unlock()
	}
}

func (c Cache) Add(url string, data []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	entry := CacheEntry{
		createdAt: time.Now(),
		val:       data,
	}
	c.Entries[url] = entry
}

func (c Cache) Get(url string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	data, ok := c.Entries[url]
	if !ok {
		return nil, false
	}
	return data.val, true
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		Entries:              make(map[string]CacheEntry),
		Mu:                   &sync.Mutex{},
		InvalidationInterval: interval,
	}
	go cache.ReapLoop()
	return cache
}
