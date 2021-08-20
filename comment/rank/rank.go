package rank

import (
	"fmt"
	"sync"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
)

type Rank struct {
	uid *sync.Map
	db  db.DB
}

func NewRank(db_ db.DB) *Rank {
	uid := &sync.Map{}
	r := &Rank{
		db:  db_,
		uid: uid,
	}

	users, err := db_.Find(&entry.User{}, db.Param{
		Order: "id asc",
		Page:  -1,
	})
	if err != nil {
		return r
	}

	for _, user := range *users.(*[]entry.User) {
		uid.Store(fmt.Sprint(user.UID), fmt.Sprint(user.ID))
	}

	return r
}

func (r Rank) Do(page, size int, t, sort string, uids ...string) (interface{}, error) {
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

	for i := range uids {
		if id, ok := r.uid.Load(uids[i]); ok {
			uids[i] = id.(string)
		}
	}

	return r.db.Find(&entry.Comment{}, db.Param{
		Page:  page,
		Size:  size,
		Order: sort,
		Query: "time > ? and user_id in (?)",
		Args:  []interface{}{t, uids},
	})
}
