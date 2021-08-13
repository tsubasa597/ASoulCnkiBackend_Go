package cache

import "github.com/tsubasa597/ASoulCnkiBackend/db/entry"

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Set(interface{}, interface{}) error
	Update(interface{}, interface{}) error
	Save() error
	Increment(entry.Comment, map[int64]struct{}) error
}
