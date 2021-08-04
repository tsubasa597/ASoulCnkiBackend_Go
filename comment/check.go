package comment

import (
	"container/heap"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

func Compare(s string) CompareResults {
	commResults := make(CompareResults, 0, conf.HeapLength)
	counts := make(map[string]float64)
	for _, v := range Hash(s) {
		val, err := commCache.Get(strconv.Itoa(int(v)))
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

		i, _ := strconv.Atoi(id)
		comm := &db.Comment{
			Model: db.Model{
				ID: uint64(i),
			},
		}
		if err := db.Find(comm); err != nil {
			continue
		}

		if utf8.RuneCountInString(comm.Comment) >= charNum {
			heap.Push(&commResults, CompareResult{
				Comment:    comm,
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
