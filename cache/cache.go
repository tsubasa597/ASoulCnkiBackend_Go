package cache

type Cacher interface {
	Get(string) (string, error)
	Set(string, string) error
	Save() error
	Increment(string, map[int64]struct{}) error
	Stop()
}

type Cache struct {
	Check   Cacher
	Content Cacher
}

func New() (*Cache, error) {
	check, err := NewBuntDB("/check")
	if err != nil {
		return nil, err
	}

	content, err := NewBuntDB("/content")
	if err != nil {
		return nil, err
	}

	return &Cache{
		Check:   *check,
		Content: *content,
	}, nil
}
