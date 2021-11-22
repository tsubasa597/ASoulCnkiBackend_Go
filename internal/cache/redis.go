package cache

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

// Redis 缓存
type Redis struct {
	db  *redis.Client
	ctx context.Context
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
		db:  db,
		ctx: ctx,
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
func (r Redis) Increment(key string, field string, value interface{}) error {
	switch val := value.(type) {
	case map[int64]struct{}:
		for k := range val {
			res := map[int64]struct{}{k: {}}

			b, err := r.db.HGet(r.ctx, key, strconv.Itoa(int(k))).Bytes()
			if err != nil {
				v, err := Serialize(res)
				if err != nil {
					continue
				}

				r.db.HSet(r.ctx, key, strconv.Itoa(int(k)), v)
				continue
			}

			if err := Deserialize(b, &res); err != nil {
				continue
			}

			v, err := Serialize(res)
			if err != nil {
				continue
			}

			r.db.HSet(r.ctx, key, strconv.Itoa(int(k)), v)
		}
	case string:
		v, err := Serialize(val)
		if err != nil {
			return err
		}

		r.db.HSet(r.ctx, key, field, v)
	}

	return nil
}

// Stop 停止
func (r Redis) Stop() error {
	// r.db.ShutdownSave(ctx)
	return nil
}
