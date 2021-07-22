package cache

type Cache interface {
	Load(interface{}) (interface{}, bool)
	Store(interface{}, interface{})
	Range(func(interface{}, interface{}) bool)
	New()
}
