package localcache

import (
	"sync"
	"time"
)

var (
	timeNow = time.Now
	ttl = 30 * time.Minute
	checkInterval = 1 * time.Second
)

type cache struct {
	store map[string]value
	m sync.Mutex
	stop chan bool
}

func New() *cache {
	c := &cache{
		store: map[string]value{},
		stop: make(chan bool),
	}
	go c.deleteExpiredKey()

	return c
}

// Get cache value by key
func (c *cache) Get(key string) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()
	v, ok := c.store[key]

	if !ok {
		return nil, NewKeyNotFound(key)
	}

	if v.isExpired() {
		_ = c.delete(key)
		return nil, NewKeyNotFound(key)
	}

	return v.v, nil
}

// Set value with associate key
func (c *cache) Set(k string, v interface{}) error {
	expireAt := time.Time.Add(timeNow(), ttl)
	c.m.Lock()
	defer c.m.Unlock()
	c.store[k] = value{
		v:         v,
		expiredAt: expireAt,
	}
	return nil
}

func (c *cache) Stop() {
	c.stop <- true
}

func (c *cache) delete(key string) error {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.store, key)
	return nil
}

func (c *cache) deleteExpiredKey() {
	for {
		select {
			case <- c.stop:
				return
		default:
			for k, v := range c.store {
				if !v.isExpired() {
					continue
				}
				_ = c.delete(k)
			}
			time.Sleep(checkInterval)
		}
	}
}

type value struct {
	v interface{}
	expiredAt time.Time
}

func (v value) isExpired() bool {
	now := timeNow()
	return v.expiredAt.Before(now)
}
