package sredis

import (
	"github.com/gomodule/redigo/redis"
	"strings"
	"time"
)

type RedisClient struct {
	pool *redis.Pool // redis pool can single mode or other mode
}

func NewRedisClient(redisConf MyRedisConf) (*RedisClient, error) {
	pool, err := NewRedis(&redisConf)
	if err != nil {
		return nil, err
	}
	return &RedisClient{pool: pool}, nil
}

// MyRedisConf redis config
type MyRedisConf struct {
	RedisHost string `yaml:"host"`

	// Maximum number of idle connections in the pool.
	RedisMaxIdle int `yaml:"max_idle"`

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	RedisMaxActive int `yaml:"max_active"`

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	RedisIdleTimeout int    `yaml:"idle_timeout"`
	RedisDB          int    `yaml:"database"`
	RedisPass        string `yaml:"pass"`
	IsCluster        bool   `yaml:"is_cluster"`  // sentinel
	MasterName       string `yaml:"master_name"` // sentinel
}

// NewRedis new a redis pool
func NewRedis(redisConf *MyRedisConf) (pool *redis.Pool, err error) {
	// sentinel use other func
	if redisConf.IsCluster {
		return InitSentinelRedisPool(redisConf)
	}
	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.Dial("tcp", redisConf.RedisHost,
				redis.DialPassword(redisConf.RedisPass),
				redis.DialDatabase(redisConf.RedisDB),
				redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout), redis.DialWriteTimeout(timeout))
			if err != nil {
				return c, err
			}
			return c, nil
		},
	}

	conn := pool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	_, err = conn.Do("ping")
	return
}

func InitSentinelRedisPool(redisConf *MyRedisConf) (pool *redis.Pool, err error) {
	s := &Sentinel{
		Address:    strings.Split(redisConf.RedisHost, ","),
		MasterName: redisConf.MasterName,
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 1000 * time.Millisecond
			c, err := redis.Dial("tcp", addr, redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout), redis.DialWriteTimeout(timeout))
			if err != nil {
				return c, err
			}
			return c, nil
		},
	}

	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (c redis.Conn, err error) {
			masterAddr, err := s.MasterAddress()
			if err != nil {
				return
			}

			timeout := 1000 * time.Millisecond

			c, err = redis.Dial("tcp", masterAddr, redis.DialPassword(redisConf.RedisPass), redis.DialConnectTimeout(timeout),
				redis.DialReadTimeout(timeout), redis.DialWriteTimeout(timeout))
			if err != nil {
				return c, err
			}

			_, err = c.Do("SELECT", redisConf.RedisDB)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}

	conn := pool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	_, err = conn.Do("ping")
	return
}
