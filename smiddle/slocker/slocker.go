package slocker

import (
	"context"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/hunterhug/gorlock"
)

var (
	// LockFactory 分布式锁工厂
	LockFactory gorlock.LockFactory
)

func Init(pool *redis.Pool) {
	LockFactory = gorlock.New(pool)
}

// Lock 获取锁
func Lock(key string, expireMillSecond int) (lock *gorlock.Lock, err error) {
	if LockFactory == nil {
		err = errors.New("redis empty")
		return
	}

	lock, err = LockFactory.Lock(context.Background(), key, expireMillSecond)
	return
}

// LockForceNotKeepAlive 强制不续命锁
func LockForceNotKeepAlive(key string, expireMillSecond int) (lock *gorlock.Lock, err error) {
	if LockFactory == nil {
		err = errors.New("redis empty")
		return
	}

	lock, err = LockFactory.LockForceNotKeepAlive(context.Background(), key, expireMillSecond)
	return
}

// UnLock 解锁
func UnLock(lock *gorlock.Lock) (success bool, err error) {
	if LockFactory == nil {
		err = errors.New("redis empty")
		return
	}

	success, err = LockFactory.UnLock(context.Background(), lock)
	return
}
