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
}

var (
	_   Cacher          = (*Redis)(nil)
	ctx context.Context = context.Background()
)

// Get 获取缓存值
func (r Redis) Get(key, field string) (string, error) {
	return r.db.HGet(ctx, key, field).Result()
}

// NewRedis 实例化 Redis
func NewRedis() (*Redis, error) {
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
	}, nil
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
			c := r.db.HGet(ctx, key, fmt.Sprint(k))
			if c.Err() == nil {
				if strings.Contains(c.Val(), field) {
					continue
				}

				r.db.HSet(ctx, key, fmt.Sprint(k), c.Val()+","+field)
				continue
			}
			r.db.HSet(ctx, key, fmt.Sprint(k), field)
		}
	case string:
		r.db.HSet(ctx, key, field, v)
	}

	return nil
}

// Stop 停止
func (r Redis) Stop() error {
	// r.db.ShutdownSave(ctx)
	return nil
}
