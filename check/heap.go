package check

import (
	"container/heap"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
)

type CompareResult struct {
	ID         uint64
	Similarity float64
}

type CompareResults []CompareResult

var _ heap.Interface = (*CompareResults)(nil)

func (r CompareResults) Len() int {
	return len(r)
}

func (r CompareResults) Less(i, j int) bool {
	if r[i].Similarity == r[j].Similarity {
		return r[i].ID < r[j].ID
	}
	return r[i].Similarity > r[j].Similarity
}

func (r CompareResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r *CompareResults) Pop() interface{} {
	old := *r
	n := len(old)

	*r = old[:n-1]

	return old[n-1]
}

func (r *CompareResults) Push(data interface{}) {
	if len(*r) > conf.HeapLength {
		*r = (*r)[:conf.HeapLength]
	}
	*r = append(*r, data.(CompareResult))
}
