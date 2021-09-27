package cache

import "sync"

type Cacher interface {
	Get(string, string) (string, error)
	Save() error
	Increment(string, string, interface{}) error
	Stop()
}

var (
	cache Cache
	once  sync.Once = sync.Once{}
)

type Cache struct {
	Check   Cacher
	Content Cacher
}

func (c Cache) Stop() {
	c.Check.Save()
	c.Content.Save()
	c.Check.Stop()
	c.Content.Stop()
}

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

func GetCache() Cache {
	return cache
}
