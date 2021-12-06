package cache

import (
	"admin/config"
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

var cacheType string

func SetupCache(cfg *config.Config) error {
	cacheType = cfg.Cache
	if cfg.Cache == "redis" {
		return setupRedisCache(cfg.Redis)
	} else if cfg.Cache == "memory" {
		return setupMemoryCache(cfg.Memory)
	}

	return fmt.Errorf("%s", "invalid cache type")
}

func Get(key string) (interface{}, error) {
	if cacheType == "" {
		return "", fmt.Errorf("%s", "invalid cache type")
	}
	if cacheType == "redis" {
		return redisGet(key)
	} else if cacheType == "memory" {
		return memoryGet(key)
	}

	return "", nil
}

func Set(key string, value interface{}, ttl time.Duration) error {
	if cacheType == "" {
		return fmt.Errorf("%s", "invalid cache type")
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(value)
	if err != nil {
		return err
	}

	if cacheType == "redis" {
		return redisSet(key, buffer.Bytes(), ttl)
	} else if cacheType == "memory" {
		return memorySet(key, buffer.Bytes(), ttl)
	}

	return nil
}
