package cache

import (
	"admin/config"
	"github.com/patrickmn/go-cache"
	"time"
)

var memoryConfig *config.Memory
var memoryCache *cache.Cache

func setupMemoryCache(cfg *config.Memory) error {
	memoryConfig = cfg
	memoryCache = cache.New(5*time.Minute, memoryConfig.PurgeTime)

	return nil
}

func memoryGet(key string) (interface{}, error) {
	if x, found := memoryCache.Get(key); found {
		return x, nil
	}
	return nil, nil
}

func memorySet(key string, value interface{}, expire time.Duration) error {
	memoryCache.Set(key, value, expire)
	return nil
}
