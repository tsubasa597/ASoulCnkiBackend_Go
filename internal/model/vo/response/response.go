package response

import "sync"

const (
	// TextLong 查重字符串长度过短
	TextLong = 20001
	// TextShort 查重字符串长度过长
	TextShort = 20002
	// ParamErr 参数错误
	ParamErr = 20003
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
	StartTime int64   `json:"start_time"`
	EndTme    int64   `json:"end_time"`
	AllCount  int64   `json:"all_count"`
}

// Related 查重
type Relateds struct {
	Rate      float64   `json:"rate"`
	Related   []Related `json:"related"`
	StartTime int64     `json:"start_time"`
	EndTme    int64     `json:"end_time"`
}

type Related struct {
	Rate  float64 `json:"rate"`
	Reply Reply   `json:"reply"`
	URL   string  `json:"reply_url"`
}

var (
	responsePool sync.Pool = sync.Pool{
		New: func() interface{} {
			return &Response{}
		},
	}
)

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
	case TextLong:
		r.Code = TextLong
		r.Message = "text to check too long"
	case TextShort:
		r.Code = TextShort
		r.Message = "text to check too short"
	case ParamErr:
		r.Code = ParamErr
		r.Message = "param error"
	default:
		r.Code = -1
		r.Message = "fail"
	}

	return r
}
