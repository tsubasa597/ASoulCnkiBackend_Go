package rank

import (
	"fmt"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/models/vo"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/model"
)

func Do(page, size int, t, sort string, uids ...string) *vo.Response {
	switch t {
	case "1":
		t = fmt.Sprint(time.Now().AddDate(0, 0, -7).Unix())
	case "2":
		t = fmt.Sprint(time.Now().AddDate(0, 0, -3).Unix())
	default:
		t = "0"
	}

	switch sort {
	case "1":
		sort = "total_like"
	case "2":
		sort = "like"
	case "3":
		sort = "num"
	}

	replys, count, err := model.Rank(page, size, t, sort, uids...)
	if err != nil {
		return nil
	}

	return vo.Sucess(&vo.Replies{
		Replies:  replys,
		AllCount: count,
	})
}
