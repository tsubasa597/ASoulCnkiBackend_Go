package check

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
)

type Check struct {
	db    db.DB
	cache cache.Cacher
}

func New(db db.DB, cache cache.Cacher) Check {
	return Check{
		db:    db,
		cache: cache,
	}
}

func (check Check) Compare(s string) CompareResults {
	commResults := make(CompareResults, 0, conf.HeapLength)
	counts := make(map[string]float64)
	for _, v := range Hash(s) {
		val, err := check.cache.Get(fmt.Sprint(v))
		if err != nil {
			continue
		}
		for _, id := range strings.Split(val.(string), ",") {
			if len(id) == 0 {
				continue
			}
			counts[id] += 1.0
		}
	}

	for id, count := range counts {
		charNum := utf8.RuneCountInString(s)
		id_, _ := strconv.Atoi(id)
		comms, err := check.db.Find(&entry.Comment{}, db.Param{
			Query: "id = ?",
			Args:  []interface{}{id_},
			Order: "id asc",
		})
		if err != nil {
			continue
		}

		comm := (*comms.(*[]entry.Comment))[0]
		if utf8.RuneCountInString(comm.Comment) >= charNum {
			heap.Push(&commResults, CompareResult{
				Comment:    &comm,
				Similarity: count / float64(utf8.RuneCountInString(comm.Comment)-conf.DefaultK+1),
			})
		}
	}

	for i, v := range commResults {
		if v.Similarity < 0.2 {
			return commResults[:i]
		}
	}

	return commResults
}

func ReplaceStr(s string) string {
	return replacer.Replace(s)

}

var (
	replacer = strings.NewReplacer("\n", "", " ", "")
)