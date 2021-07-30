package comment

import (
	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

var (
	caches = &cache.Cache{}
)

func InitCache() {
	caches.Once.Do(func() {
		caches = cache.Init()
	})
}

func LoadCache() {
	for _, v := range *db.Get(db.Comment{}).(*[]db.Comment) {
		c := v
		for k := range HashSet(v.Comment) {
			if v, ok := caches.Comm.Load(k); ok {
				if comma, ok := v.(*Set); ok {
					(*comma)[int64(c.ID)] = struct{}{}
				}
			} else {
				s := make(Set)
				s[int64(c.ID)] = struct{}{}
				caches.Comm.Store(k, &s)
			}
		}
	}
}
