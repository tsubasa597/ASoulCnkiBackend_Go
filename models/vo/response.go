package vo

import "sync"

const (
	// TEXTTOLONG 查重字符串长度过短
	TEXTTOLONG = 20001
	// TEXTTOSHORT 查重字符串长度过长
	TEXTTOSHORT = 20002
)

// Response 返回数据结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Replies 作文展
type Replies struct {
	Replies   []Reply `json:"replies"`
	StartTime int32   `json:"start_time"`
	EndTme    int32   `json:"end_time"`
	AllCount  int64   `json:"all_count"`
}

// Related 查重
type Related struct {
	Rate      float64 `json:"rate"`
	URL       string  `json:"reply_url"`
	Related   []Reply `json:"related"`
	StartTime int32   `json:"start_time"`
	EndTme    int32   `json:"end_time"`
	AllCount  int     `json:"all_count"`
}

// Sucess 请求成功
func Sucess(data interface{}) *Response {
	r := responsePool.Get().(*Response)
	defer responsePool.Put(r)

	r.Code = 0
	r.Message = "sucess"
	r.Data = data
	return r
}

// Fail 请求失败
func Fail(code int) *Response {
	r := responsePool.Get().(*Response)
	defer responsePool.Put(r)

	r.Data = nil
	switch code {
	case TEXTTOLONG:
		r.Code = TEXTTOLONG
		r.Message = "text to check too long"
	case TEXTTOSHORT:
		r.Code = TEXTTOSHORT
		r.Message = "text to check too short"
	default:
		r.Code = -1
		r.Message = "fail"
	}

	return r
}

var (
	responsePool sync.Pool = sync.Pool{
		New: func() interface{} {
			return &Response{}
		},
	}
)
