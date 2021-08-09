package cache

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Set(interface{}, interface{}) error
	Save() error
	Init(interface{}, interface{}) error
}
