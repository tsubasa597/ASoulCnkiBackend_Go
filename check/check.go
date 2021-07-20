package check

import (
	"container/heap"
	"sync"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

var (
	data   sync.Map
	length int
)

func Init() {
	for _, v := range *db.Get(db.Comment{}).(*[]db.Comment) {
		data.Store(v.ID, HashSet(v.Comment))
		length++
	}
}

type Result struct {
	Comment    db.Comment
	Similarity float64
}

func Compare(s string) []Result {
	h1 := Hash(s)
	comResults := make(CompareResults, 0, length)
	data.Range(func(key, value interface{}) bool {
		set := make(Set)
		count := 0.0
		charNum := utf8.RuneCountInString(s)
		for i := 0; i < charNum-conf.DefaultK+1; i++ {
			if _, ok := value.(Set)[h1[i]]; ok {
				for j := 0; j < conf.DefaultK; j++ {
					set[int64(i+j)] = struct{}{}
				}
			}
		}
		for i := 0; i < charNum; i++ {
			if _, ok := set[int64(i)]; ok {
				count++
			}
		}

		heap.Push(&comResults, CompareResult{
			ID:         key.(uint64),
			Similarity: count / float64(charNum),
		})
		return true
	})

	result := make([]Result, 0, conf.HeapLength)
	for _, v := range comResults {
		if v.Similarity < 0.2 {
			continue
		}

		comm := &db.Comment{
			Model: db.Model{
				ID: v.ID,
			},
		}

		db.Find(comm)
		result = append(result, Result{
			Comment:    *comm,
			Similarity: v.Similarity,
		})
	}

	return result
}
