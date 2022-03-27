package cache

import (
	"pmain2/internal/apperror"
	"sync"
	"time"
)

type CacheI interface {
	Get(key, t interface{}) (interface{}, bool)
	Set(key, value interface{}, ttl time.Duration)
}

type Item struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	item              map[interface{}]Item
}

func CreateCache(dExpiration, cInterval time.Duration) *Cache {
	item := make(map[interface{}]Item)

	c := Cache{
		item:              item,
		defaultExpiration: dExpiration,
		cleanupInterval:   cInterval,
	}

	c.StartGC()

	return &c
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	item, ok := c.item[key]
	if !ok {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Set(key, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	expiration := ttl
	if expiration == 0 {
		expiration = c.defaultExpiration
	}

	c.item[key] = Item{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *Cache) Delete(key interface{}) error {
	c.Lock()

	defer c.Unlock()

	if _, found := c.item[key]; !found {
		return apperror.ErrCacheKeyNotFound
	}

	delete(c.item, key)

	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.item == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) expiredKeys() (keys []interface{}) {
	c.RLock()

	defer c.RUnlock()

	for k, i := range c.item {
		if time.Now().UnixNano() > i.Expiration {
			keys = append(keys, k)
		}
	}

	return
}

func (c *Cache) clearItems(keys []interface{}) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.item, k)
	}
}
