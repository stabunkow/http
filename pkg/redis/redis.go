package redis

import (
	"stabunkow/http/pkg/setting"
	"time"

	"github.com/gomodule/redigo/redis"
)

var DefaultCache *Cache

func GetDefaultCache() *Cache {
	return DefaultCache
}

func Setup() error {
	DefaultCache = NewCache(setting.RedisSetting.Addr)
	return nil
}

type Cache struct {
	*redis.Pool
}

func NewCache(addr string) *Cache {
	pool := &redis.Pool{
		MaxIdle:     30,
		MaxActive:   30,
		IdleTimeout: 200,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &Cache{pool}
}

func (c *Cache) Conn() redis.Conn {
	return c.Get()
}
