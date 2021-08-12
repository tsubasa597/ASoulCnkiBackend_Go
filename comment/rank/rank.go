package rank

import (
	"fmt"

	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
)

type Rank struct {
	db db.DB
}

func NewRank(db db.DB) *Rank {
	return &Rank{
		db: db,
	}
}

func (r Rank) R() {
	fmt.Println(r.db.Find(&entry.Article{}, db.NewFilter(db.All, []entry.User{}, db.TotalLikeSort)))

}
