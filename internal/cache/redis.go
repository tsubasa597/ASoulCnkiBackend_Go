package cache

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

// Redis 缓存
type Redis struct {
	db    *redis.Client
	mutex *sync.Mutex
	ctx   context.Context
}

var (
	_ Cacher = (*Redis)(nil)
)

// NewRedis 实例化 Redis
func NewRedis() (*Redis, error) {
	ctx := context.Background()

	db := redis.NewClient(&redis.Options{
		Addr:     config.RedisADDR,
		Password: config.RedisPwd,
		DB:       config.DB,
	})
	if _, err := db.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Redis{
		db:    db,
		mutex: &sync.Mutex{},
		ctx:   ctx,
	}, nil
}

// Get 获取缓存值
func (r Redis) Get(key, field string) (string, error) {
	return r.db.HGet(r.ctx, key, field).Result()
}

// Save 持久化
func (r Redis) Save() error {
	return nil
}

// Increment 添加数据
func (r Redis) Increment(key string, field string, val interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	switch v := val.(type) {
	case map[int64]struct{}:
		for k := range v {
			c := r.db.HGet(r.ctx, key, fmt.Sprint(k))
			if c.Err() == nil {
				if strings.Contains(c.Val(), field) {
					continue
				}

				r.db.HSet(r.ctx, key, fmt.Sprint(k), c.Val()+","+field)
				continue
			}
			r.db.HSet(r.ctx, key, fmt.Sprint(k), field)
		}
	case string:
		r.db.HSet(r.ctx, key, field, v)
	}

	return nil
}

// Stop 停止
func (r Redis) Stop() error {
	// r.db.ShutdownSave(ctx)
	return nil
}
