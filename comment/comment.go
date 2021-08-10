package comment

import (
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/comment/check"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
)

type Comment struct {
	*ListenUpdate
	check.Check
}

func New(db_ db.DB, cache cache.Cacher, log *logrus.Entry) *Comment {
	c := &Comment{
		ListenUpdate: NewListen(db_, cache, log),
		Check:        check.New(db_, cache),
	}

	if c.ListenUpdate.enable {
		users, err := c.db.Find(&entry.User{}, db.Param{
			Order: "id asc",
		})
		if err != nil {
			return c
		}

		for _, user := range *users.(*[]entry.User) {
			c.ListenUpdate.Add(user)
		}
	}

	if err := cache.Increment(db_, check.HashSet); err != nil {
		c.log.WithField("Func", "cache.Increment").Error(err)
	}

	return c
}
