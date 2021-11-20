package request

// Check 请求 api/v1/check 接口数据结构体
type Check struct {
	Text string `form:"text"`
}
