package rank

import (
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

type Rank struct {
	db db.DB
}

func NewRank(db db.DB) *Rank {
	return &Rank{
		db: db,
	}
}

func (r Rank) Do(time, sort string, uids ...string) (interface{}, error) {
	return r.db.Rank(time, sort, uids), nil
	// return r.db.Find(&entry.Comments{}, param)
}
