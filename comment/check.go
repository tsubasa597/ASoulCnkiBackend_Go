package comment

import (
	"container/heap"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

func ReplaceStr(s string) string {
	emoteCache.Range(func(key, value interface{}) bool {
		s = strings.Replace(s, key.(string), value.(string), -1)
		return true
	})
	return s
}

func ReplaceRune(s string) string {
	emoteCache.Range(func(key, value interface{}) bool {
		s = strings.Replace(s, value.(string), key.(string), -1)
		return true
	})
	return s
}

func Compare(s string) CompareResults {
	h1 := Hash(s)

	commResults := make(CompareResults, 0, conf.HeapLength)
	commCache.Range(func(key, value interface{}) bool {
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

		heap.Push(&commResults, CompareResult{
			Comm:       key.(*db.Comment),
			Similarity: count / float64(charNum),
		})
		return true
	})

	for i, v := range commResults {
		if v.Similarity < 0.2 {
			return commResults[:i]
		}

		v.Comm.Comment = ReplaceRune(v.Comm.Comment)
	}

	return commResults
}
