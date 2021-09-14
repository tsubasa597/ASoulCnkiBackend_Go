package comment

import (
	"fmt"
	"time"

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

func New(dD db.DB, cache cache.Cache, log *logrus.Entry) *Comment {
	c := &Comment{
		Rank:         *rank.NewRank(dD),
		ListenUpdate: update.NewListen(dD, cache, log),
		Check:        check.New(dD, cache),
	}

	if c.ListenUpdate.Enable {
		users, err := dD.Find(&entry.User{}, db.Param{
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

	val, err := cache.Content.Get("LastCommentID")
	if err != nil {
		val = "0"
	}
	comms := dD.GetContent(val)

	for _, comm := range comms {
		if err := cache.Content.Set(fmt.Sprint(comm.Rpid), comm.Content); err != nil {
			log.WithField("Func", "cache.Set").Error(err)
		}
		cache.Content.Set("LastCommentID", fmt.Sprint(comm.Rpid))
	}

	if err := cache.Content.Save(); err != nil {
		log.WithField("Func", "cache.Save").Error(err)
	}

	val, err = cache.Check.Get("LastCommentID")
	if err != nil {
		val = "0"
	}
	comms = dD.GetContent(val)

	for _, comm := range comms {
		if err := cache.Check.Increment(fmt.Sprint(comm.Rpid), check.HashSet(comm.Content)); err != nil {
			log.WithField("Func", "cache.Increment").Error(err)
		}
	}

	if err := cache.Check.Save(); err != nil {
		log.WithField("Func", "cache.Save").Error(err)
	}

	return c
}

var (
	instance *Comment
)

func GetInstance() *Comment {
	return instance
}

func init() {
	db, err := db.New()
	if err != nil {
		panic(err)
	}

	cache, err := cache.New()
	if err != nil {
		panic(err)
	}

	entry := logrus.NewEntry(logrus.StandardLogger()).
		WithField("Time", time.Now().Format("2006/01/02 15:04:05"))

	instance = New(*db, *cache, entry)
}
