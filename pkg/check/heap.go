package check

import (
	"container/heap"

	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
)

// CompareResult 查重结果
type CompareResult struct {
	ID         string
	Similarity float64
}

// CompareResults 查重结果切片
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

// Pop 弹出
func (r *CompareResults) Pop() interface{} {
	old := *r
	*r = old[1:]

	return old[0]
}

// Push 推入
func (r *CompareResults) Push(data interface{}) {
	if len(*r) == setting.HeapLength {
		*r = (*r)[:setting.HeapLength-1]
	}
	*r = append(*r, data.(CompareResult))
}
