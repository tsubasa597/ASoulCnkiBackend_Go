package request

// Rank 请求 api/v1/rank 接口结构体
type Rank struct {
	Page int    `form:"page"`
	Size int    `form:"size"`
	Time int64  `form:"time"`
	Sort string `form:"sort"`
	Ids  []int  `form:"ids"`
}
