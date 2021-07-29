package cache

import "sync"

type Reply struct {
	data *sync.Map
}

var _ Cacher = (*Reply)(nil)

func (r Reply) Load(v interface{}) (interface{}, bool) {
	return r.data.Load(v)
}

func (r Reply) Store(key, value interface{}) {
	r.data.Store(key, value)
}

func (r Reply) Range(f func(interface{}, interface{}) bool) {
	r.data.Range(f)
}

func (r *Reply) New() {
	r.data = &sync.Map{}
}
