package cache

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
)

type Redis struct {
	db    *redis.Client
	mutex *sync.Mutex
}

var (
	_   Cacher          = (*Redis)(nil)
	ctx context.Context = context.Background()
)

func (r Redis) Get(key, field string) (string, error) {
	return r.db.HGet(ctx, key, field).Result()
}

func NewRedis() (*Redis, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     setting.RedisADDR,
		Password: setting.RedisPwd,
		DB:       setting.DB,
	})
	if _, err := db.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Redis{
		db:    db,
		mutex: &sync.Mutex{},
	}, nil
}

func (r Redis) Save() error {
	return nil
}

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
		r.db.HSet(ctx, key, "LastCommentID", field)
	case string:
		r.db.HSet(ctx, key, field, v)
		r.db.HSet(ctx, key, "LastCommentID", v)
	}

	return nil
}

func (r Redis) Stop() {
	// r.db.ShutdownSave(ctx)
}
