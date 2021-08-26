package check

import (
	"container/heap"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/vo"
)

type Check struct {
	db    db.DB
	cache cache.Cache
}

func New(db db.DB, cache cache.Cache) Check {
	return Check{
		db:    db,
		cache: cache,
	}
}

func (check Check) Compare(s string) vo.Related {
	commResults := make(CompareResults, 0, conf.HeapLength)
	counts := make(map[string]float64)
	for _, v := range Hash(s) {
		val, err := check.cache.Check.Get(fmt.Sprint(v))
		if err != nil {
			continue
		}
		for _, id := range strings.Split(val, ",") {
			if len(id) == 0 {
				continue
			}
			counts[id] += 1.0
		}
	}

	for id, count := range counts {
		charNum := utf8.RuneCountInString(s)

		content, err := check.cache.Content.Get(id)
		if err != nil {
			continue
		}

		n := utf8.RuneCountInString(ReplaceStr(string(content)))
		if n >= charNum {
			heap.Push(&commResults, CompareResult{
				ID:         id,
				Similarity: count / float64(n-conf.DefaultK+1),
			})
		}
	}

	related := vo.Related{
		Related: make([]vo.Reply, 0),
	}
	for len(commResults) > 0 {
		comresult := commResults.Pop().(CompareResult)
		if comresult.Similarity < 0.2 {
			break
		}

		reply, err := check.db.Check(comresult.ID)
		if err != nil {
			continue
		}
		related.Rate = comresult.Similarity
		related.Related = append(related.Related, reply)
	}
	return related
}

func ReplaceStr(s string) string {
	return replacer.Replace(s)
}

var (
	replacer = strings.NewReplacer("\n", "", " ", "")
)
