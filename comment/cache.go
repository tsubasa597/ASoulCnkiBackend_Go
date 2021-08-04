package comment

import (
	"strconv"
	"strings"

	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

var (
	commCache cache.Cacher
)

func InitCache() {
	var err error
	commCache, err = cache.NewComment(func(commCache cache.Cacher) {
		for _, v := range *db.Get(db.Comment{}).(*[]db.Comment) {
			for k := range HashSet(v.Comment) {
				data, id := strconv.Itoa(int(k)), strconv.Itoa(int(v.ID))
				if ids, err := commCache.Get(data); err == nil {
					if strings.Contains(ids.(string), id) {
						continue
					}
					commCache.Set(data, ids.(string)+","+id)
					continue
				}
				commCache.Set(data, id)
			}
		}
	})
	if err != nil {
		panic(err)
	}
}
