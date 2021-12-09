package cache

import (
	"admin/config"
	"fmt"
	"time"
)

const (
	MEMORY = 1
	REDIS  = 2
)

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
}

var cache map[int]interface{}
var cacheType string

func SetupCache(cfg *config.Config) error {
	cache := make(map[int]interface{})
	cacheType = cfg.Cache

	cacheRedis, err := setupRedis(cfg.Redis)
	if err != nil {
		return err
	}
	cache[REDIS] = cacheRedis

	memory, err := setupMemory(cfg.Memory)
	if err != nil {
		return err
	}
	cache[MEMORY] = memory

	return nil
}

func Get(cacheType int, key string) (interface{}, error) {
	instance, ok := cache[cacheType]
	if ok {
		return instance.(Cache).Get(key)
	}

	return nil, fmt.Errorf("%s", "unknown cache type")
}

func Set(cacheType int, key string, value interface{}, ttl time.Duration) error {
	instance, ok := cache[cacheType]
	if ok {
		return instance.(Cache).Set(key, value, ttl)
	}

	return fmt.Errorf("%s", "unknown cache type")
}
