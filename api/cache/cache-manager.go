package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var Sessions *Cache

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Value      *Session
	Created    time.Time
	Expiration int64
}

type Session struct {
	SessionID string
	Belongs   int64
}

func NewManager(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	if cleanupInterval > 0 {
		cache.StartGC()
	}
	return &cache
}

func (c *Cache) Set(key string, value *Session, duration time.Duration) {
	var expiration int64
	if duration == 0 {
		duration = c.defaultExpiration
	}
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.Lock()
	defer c.Unlock()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
	fmt.Println(key)
}

func (c *Cache) Get(key string) (*Session, bool) {
	c.RLock()
	defer c.RUnlock()
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Value, true
}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}
	delete(c.items, key)
	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}
func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)
		if c.items == nil {
			return
		}
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()
	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}

func (c *Cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()
	for _, k := range keys {
		delete(c.items, k)
	}
}
