package comment

import (
	"sync"

	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

var (
	once   = sync.Once{}
	caches = &cache.Cache{}
)

func InitCache() {
	once.Do(func() {
		caches = cache.Init()
	})
}

func LoadCache() {
	for _, v := range *db.Get(db.Comment{}).(*[]db.Comment) {
		c := v
		for k := range HashSet(v.Comment) {
			if v, ok := caches.Comm.Load(k); ok {
				if comma, ok := v.(*Set); ok {
					(*comma)[c.UID] = struct{}{}
				}
			} else {
				s := make(Set)
				s[c.UID] = struct{}{}
				caches.Comm.Store(k, &s)
			}
		}
		caches.Reply.Store(c.UID, &c)
	}

	for _, v := range *db.Get(db.Emote{}).(*[]db.Emote) {
		caches.Reply.Store(v.EmoteText, string(v.EmoteID))
	}
}
