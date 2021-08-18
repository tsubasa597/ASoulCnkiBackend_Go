package cache

type Cacher interface {
	Get(string) (string, error)
	Set(string, string) error
	Save() error
}

type Cache struct {
	Check   LevelDB
	Content LevelDB
}

func New() (*Cache, error) {
	check, err := NewLevelDB("/check")
	if err != nil {
		return nil, err
	}

	content, err := NewLevelDB("/content")
	if err != nil {
		return nil, err
	}
	return &Cache{
		Check:   *check,
		Content: *content,
	}, nil
}
