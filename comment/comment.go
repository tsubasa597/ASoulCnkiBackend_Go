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

	id, err := c.cache.Get("LastCommentID")
	if err != nil {
		for _, v := range *c.db.Get(entry.Comment{}).(*[]entry.Comment) {
			if err := c.cache.Init(v.ID, check.HashSet(v.Comment)); err != nil {
				log.WithField("Func", "Init Cache").Error(err)
				return c
			}
			c.cache.Set("LastCommentID", fmt.Sprint(v.ID))
		}
	}

	comms, err := c.db.Find(&entry.Comment{}, db.Param{
		Where: map[string]interface{}{
			"ID": id.(string),
		},
		Query: "id",
		Args:  []interface{}{id},
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

	return c
}
