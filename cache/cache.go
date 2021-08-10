package cache

import (
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Set(interface{}, interface{}) error
	Update(interface{}, interface{}) error
	Save() error
	Init(interface{}, interface{}) error
	Increment(db.DB, func(string) map[int64]struct{}) error
}
