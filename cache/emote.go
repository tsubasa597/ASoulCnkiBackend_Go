package cache

import "sync"

type Emote struct {
	data *sync.Map
}

var _ Cache = (*Emote)(nil)

func (e Emote) Load(v interface{}) (interface{}, bool) {
	return e.data.Load(v)
}

func (e Emote) Store(key, value interface{}) {
	e.data.Store(key, value)
}

func (e Emote) Range(f func(interface{}, interface{}) bool) {
	e.data.Range(f)
}

func (e *Emote) New() {
	e.data = &sync.Map{}
}
