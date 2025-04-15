package sredis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

// Set help func to set redis key which use MULTI order
func (s *RedisClient) Set(key string, value []byte, expireSecond int64) (err error) {
	conn := s.pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send("SET", key, value)
	if err != nil {
		return err
	}

	if expireSecond >= 0 {
		err = conn.Send("EXPIRE", key, expireSecond)
		if err != nil {
			return err
		}
	}

	_, err = conn.Do("EXEC")
	return
}

func (s *RedisClient) SetNX(key string, value []byte, expireSecond int64) (set bool, err error) {
	conn := s.pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	reply, err := conn.Do("SET", key, string(value), "NX", "PX", expireSecond*1000)
	if err != nil {
		return false, err
	}

	if reply == nil {
		return false, nil
	}

	return true, nil
}

func (s *RedisClient) Release(key string, value []byte) (err error) {
	conn := s.pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end;
	`
	script := redis.NewScript(1, luaScript)
	_, err = script.Do(conn, key, string(value))
	return err
}

// Delete help func to delete redis key
func (s *RedisClient) Delete(key string) (err error) {
	conn := s.pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	_, err = conn.Do("DEL", key)
	return err
}

// Expire help func to long redis key expire time
func (s *RedisClient) Expire(key string, expireSecond int64) (err error) {
	conn := s.pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	if expireSecond <= 0 {
		expireSecond = 1
	}
	_, err = conn.Do("EXPIRE", key, expireSecond)
	return err
}

// Keys help func to keys redis
func (s *RedisClient) Keys(pattern string) (result []string, exist bool, err error) {
	conn := s.pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	keys, err := redis.ByteSlices(conn.Do("KEYS", pattern))
	if errors.Is(err, redis.ErrNil) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	result = make([]string, len(keys))
	for k, v := range keys {
		result[k] = string(v)
	}
	return result, true, nil
}

// Get help func to get redis key
func (s *RedisClient) Get(key string) (value []byte, ttl int64, exist bool, err error) {
	conn := s.pool.Get()
	if conn.Err() != nil {
		err = conn.Err()
		return
	}

	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	value, err = redis.Bytes(conn.Do("GET", key))
	if errors.Is(err, redis.ErrNil) {
		return nil, 0, false, nil
	} else if err != nil {
		return nil, 0, false, err
	}

	ttl, err = redis.Int64(conn.Do("TTL", key))
	if errors.Is(err, redis.ErrNil) {
		return nil, 0, false, nil
	} else if err != nil {
		return nil, 0, false, err
	}

	return value, ttl, true, nil
}

func (s *RedisClient) SetPool(pool *redis.Pool) {
	s.pool = pool
}
