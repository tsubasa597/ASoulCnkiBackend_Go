package comment

import (
	"fmt"

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

	id, err := c.cache.Get("LastCommentID")
	if err != nil {
		comms, err := c.db.Get(&entry.Comment{})
		if err != nil {
			return c
		}

		for _, v := range *comms.(*[]entry.Comment) {
			if err := c.cache.Init(v.ID, check.HashSet(check.ReplaceStr(v.Comment))); err != nil {
				log.WithField("Func", "Init Cache").Error(err)
				return c
			}
			c.cache.Set("LastCommentID", fmt.Sprint(v.ID))
		}

		c.cache.Save()
		return c
	}

	comms, err := c.db.Find(&entry.Comment{}, db.Param{
		Query: "id > ?",
		Args:  []interface{}{id},
		Order: "id",
	})
	if err != nil {
		return c
	}

	for _, v := range *comms.(*[]entry.Comment) {
		if err := c.cache.Init(v.ID, check.HashSet(v.Comment)); err != nil {
			log.WithField("Func", "Init Cache").Error(err)
			return c
		}
		c.cache.Set("LastCommentID", fmt.Sprint(v.ID))
	}

	c.cache.Save()
	return c
}
