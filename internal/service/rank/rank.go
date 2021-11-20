package rank

import (
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/internal/dao"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/request"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/response"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

const (
	week int64 = iota + 1
	threeDay

	likeNum string = "1"
	num     string = "2"
)

func Do(r request.Rank) *response.Response {
	if r.Page < 1 {
		r.Page = 1
	}

	// 每页展示评论数若超出范围，设置为默认值
	if r.Size < 1 || r.Size > 30 {
		r.Size = config.Size
	}

	// 默认查询全部时间
	switch r.Time {
	case week:
		r.Time = time.Now().AddDate(0, 0, -7).Unix()
	case threeDay:
		r.Time = time.Now().AddDate(0, 0, -3).Unix()
	default:
		r.Time = 0
	}

	// 默认为按总点赞数排序
	switch r.Sort {
	case likeNum:
		r.Sort = "like_num"
	case num:
		r.Sort = "num"
	default:
		r.Sort = "total_like"
	}

	replies, count, err := dao.Rank(r)
	if err != nil {
		return nil
	}

	startTime, endTime := dao.GetTimeInfo()
	return response.Sucess(response.Replies{
		Replies:   replies,
		AllCount:  count,
		StartTime: startTime,
		EndTme:    endTime,
	})
}
