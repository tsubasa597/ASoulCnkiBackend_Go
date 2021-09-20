package check

import (
	"container/heap"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/models/vo"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/check"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/model"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
)

func Compare(s string) vo.Related {
	commResults := make(check.CompareResults, 0, setting.HeapLength)
	counts := make(map[string]float64)
	for _, v := range check.Hash(s) {
		val, err := cache.GetCache().Check.Get(fmt.Sprint(v))
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

		content, err := cache.GetCache().Content.Get(id)
		if err != nil {
			continue
		}

		n := utf8.RuneCountInString(check.ReplaceStr(string(content)))
		if n >= charNum {
			heap.Push(&commResults, check.CompareResult{
				ID:         id,
				Similarity: count / float64(n-setting.DefaultK+1),
			})
		}
	}

	related := vo.Related{
		Related: make([]vo.Reply, 0),
	}
	for len(commResults) > 0 {
		comresult := commResults.Pop().(check.CompareResult)
		if comresult.Similarity < 0.2 {
			break
		}

		reply, err := model.Check(comresult.ID)
		if err != nil {
			continue
		}
		related.Rate = comresult.Similarity
		related.Related = append(related.Related, reply)
	}
	return related
}
