package comment

import (
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/comment/check"
	"github.com/tsubasa597/ASoulCnkiBackend/comment/rank"
	"github.com/tsubasa597/ASoulCnkiBackend/comment/update"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
)

type Comment struct {
	*update.ListenUpdate
	check.Check
	rank.Rank
}

func New(db_ db.DB, cache cache.Cacher, log *logrus.Entry) *Comment {
	c := &Comment{
		Rank:         *rank.NewRank(db_),
		ListenUpdate: update.NewListen(db_, cache, log),
		Check:        check.New(db_, cache),
	}

	if c.ListenUpdate.Enable {
		users, err := db_.Find(&entry.User{}, db.Param{
			Order: "id asc",
			Page:  -1,
		})
		if err != nil {
			return c
		}

		for _, user := range *users.(*[]entry.User) {
			c.ListenUpdate.Add(user)
		}
	}

	val, err := cache.Get("LastCommentID")
	if err != nil {
		val = "0"
	}

	comms, err := db_.Find(&entry.Comment{}, db.Param{
		Page:  -1,
		Query: "id > ?",
		Args:  []interface{}{val},
		Order: "id",
	})
	if err != nil {
		return c
	}

	for _, comm := range *comms.(*[]entry.Comment) {
		if err := cache.Increment(comm, check.HashSet(comm.Content)); err != nil {
			log.WithField("Func", "cache.Increment").Error(err)
		}
	}

	if err := cache.Save(); err != nil {
		log.WithField("Func", "cache.Save").Error(err)
	}

	return c
}
