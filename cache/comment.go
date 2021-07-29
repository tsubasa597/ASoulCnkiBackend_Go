package cache

import (
	"sync"
)

type Comment struct {
	data *sync.Map
}

var _ Cacher = (*Comment)(nil)

func (c Comment) Load(v interface{}) (interface{}, bool) {
	return c.data.Load(v)
}

func (c Comment) Store(key, value interface{}) {
	c.data.Store(key, value)
}

func (c Comment) Range(f func(interface{}, interface{}) bool) {
	c.data.Range(f)
}

func (c *Comment) New() {
	c.data = &sync.Map{}
}
