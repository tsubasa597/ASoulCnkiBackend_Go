package comment

import (
	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

var (
	commCache  cache.Comment
	emoteCache cache.Emote
)

func InitCache() {
	commCache.New()
	emoteCache.New()

	for _, v := range *db.Get(db.Comment{}).(*[]db.Comment) {
		c := v
		commCache.Store(&c, HashSet(v.Comment))
	}

	for _, v := range *db.Get(db.Emote{}).(*[]db.Emote) {
		emoteCache.Store(v.EmoteText, string(v.EmoteID))
	}
}
