package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Cache is the interface for common cache module
type Cache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
}

var c Cache

func init() {
	c = gocache.New(5*time.Minute, 10*time.Minute)
}

// Set set the key-value with ttl
func Set(k string, x interface{}, d time.Duration) {
	c.Set(k, x, d)
}

// Get get value of the key
func Get(k string) (interface{}, bool) {
	return c.Get(k)
}

// SetCache set the cache instance for the sdk
func SetCache(cache Cache) {
	c = cache
}
