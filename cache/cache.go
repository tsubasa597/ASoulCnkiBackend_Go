package cache

type Cacher interface {
	Load(interface{}) (interface{}, bool)
	Store(interface{}, interface{})
	Range(func(interface{}, interface{}) bool)
	New()
}

type Cache struct {
	Comm  Cacher
	Emote Cacher
	Reply Cacher
}

func Init() *Cache {
	c := &Cache{
		Comm:  &Comment{},
		Emote: &Emote{},
		Reply: &Reply{},
	}

	c.Comm.New()
	c.Emote.New()
	c.Reply.New()

	return c
}
