package check

import (
	"container/heap"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
)

type CompareResult struct {
	Comment    *entry.Comment
	ID         string
	Similarity float64
}

type CompareResults []CompareResult

var _ heap.Interface = (*CompareResults)(nil)

func (r CompareResults) Len() int {
	return len(r)
}

func (r CompareResults) Less(i, j int) bool {
	return r[i].Similarity > r[j].Similarity
}

func (r CompareResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r *CompareResults) Pop() interface{} {
	old := *r
	*r = old[1:]

	return old[0]
}

func (r *CompareResults) Push(data interface{}) {
	if len(*r) == conf.HeapLength {
		*r = (*r)[:conf.HeapLength-1]
	}
	*r = append(*r, data.(CompareResult))
}
