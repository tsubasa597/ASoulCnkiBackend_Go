package rank

import (
	"fmt"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/vo"
)

type Rank struct {
	db db.DB
}

func NewRank(db_ db.DB) *Rank {
	r := &Rank{
		db: db_,
	}

	return r
}

func (r Rank) Do(page, size int, t, sort string, uids ...string) *vo.Response {
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
		sort = "total_like desc"
	case "2":
		sort = "like desc"
	case "3":
		sort = "num desc"
	}

	replys, err := r.db.Rank(page, size, t, sort, uids...)
	if err != nil {
		return nil
	}

	return vo.Sucess(&vo.Replies{
		Replies: replys,
	})
}
