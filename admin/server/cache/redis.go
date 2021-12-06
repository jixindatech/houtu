package cache

import (
	"admin/config"
	"admin/core/log"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"time"
)

var redisConfig *config.Redis
var redisPool *redis.Pool

func setupRedisCache(cfg *config.Redis) error {
	redisConfig = cfg
	redisPool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: 240,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",
				cfg.Host+fmt.Sprintf(":%d", cfg.Port),
				redis.DialDatabase(cfg.Db),
				redis.DialReadTimeout(time.Duration(1)*time.Second),
				redis.DialWriteTimeout(time.Duration(1)*time.Second),
				redis.DialConnectTimeout(time.Duration(2)*time.Second),
			)
			if err != nil {
				return nil, err
			}
			if cfg.Password != "" {
				if _, err := c.Do("AUTH", cfg.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
	}

	return nil
}

func redisSet(key string, value interface{}, ttl time.Duration) error {
	if redisConfig == nil {
		return fmt.Errorf("%s", "invalid redis instance")
	}

	conn := redisPool.Get()
	defer conn.Close()

	if redisConfig.KeyPrefix != "" {
		key = redisConfig.KeyPrefix + key
	}

	_, err := conn.Do("SET", key, value, "EX", int(ttl))
	if err != nil {
		log.Logger.Error("redis", zap.String("err", err.Error()))
		return err
	}

	return err
}

func redisGet(key string) (interface{}, error) {
	if redisConfig == nil {
		return "", fmt.Errorf("%s", "invalid redis instance")
	}

	conn := redisPool.Get()
	defer conn.Close()

	if redisConfig.KeyPrefix != "" {
		key = redisConfig.KeyPrefix + key
	}

	return redis.Bytes(conn.Do("GET", key))
}
