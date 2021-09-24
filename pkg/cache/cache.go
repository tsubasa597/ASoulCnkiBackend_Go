package cache

import "sync"

type Cacher interface {
	Get(string) (string, error)
	Set(string, string) error
	Save() error
	Increment(string, map[int64]struct{}) error
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
		check, err := NewBuntDB("/check")
		if err != nil {
			panic(err)
		}

		content, err := NewBuntDB("/content")
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