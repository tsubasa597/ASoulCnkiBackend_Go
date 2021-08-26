package vo

const (
	TEXTTOLONG  = 20001
	TEXTTOSHORT = 20002
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Replies struct {
	Replies   []Reply `json:"replies"`
	StartTime int32   `json:"start_time"`
	EndTme    int32   `json:"end_time"`
	AllCount  int     `json:"all_count"`
}

type Related struct {
	Rate      float64 `json:"rate"`
	URL       string  `json:"reply_url"`
	Related   []Reply `json:"related"`
	StartTime int32   `json:"start_time"`
	EndTme    int32   `json:"end_time"`
	AllCount  int     `json:"all_count"`
}

func Sucess(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "sucess",
		Data:    data,
	}
}

func Fail(code int) *Response {
	r := &Response{}

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
