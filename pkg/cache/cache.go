package cache

import "sync"

// Cacher 缓存接口
type Cacher interface {
	Get(string, string) (string, error)
	Save() error
	Increment(string, string, interface{}) error
	Stop()
}

// Cache 缓存实例
type Cache struct {
	Check   Cacher
	Content Cacher
}

// Stop 停止
func (c Cache) Stop() {
	c.Check.Save()
	c.Content.Save()
	c.Check.Stop()
	c.Content.Stop()
}

// Setup 初始化
func Setup() {
	once.Do(func() {
		check, err := NewRedis()
		if err != nil {
			panic(err)
		}

		content, err := NewRedis()
		if err != nil {
			panic(err)
		}

		cache = Cache{
			Check:   *check,
			Content: *content,
		}
	})
}

// GetCache 获取缓存实例
func GetCache() Cache {
	return cache
}

var (
	cache Cache
	once  sync.Once = sync.Once{}
)
