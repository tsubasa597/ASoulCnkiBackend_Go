package update

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
	"github.com/tsubasa597/BILIBILI-HELPER/listen"
)

type ListenUpdate struct {
	Update
	Enable bool
	db     db.DB
	cache  cache.Cache
	log    *logrus.Entry
}

func (lis ListenUpdate) Started() bool {
	return lis.Enable
}

func (lis ListenUpdate) Stop() {
	lis.dynamic.Stop()
	lis.comment.Stop()
	lis.Enable = false
}

func (lis ListenUpdate) Add(user entry.User) {
	if !lis.Started() {
		return
	}

	ctx, ch, err := lis.Update.dynamic.Add(user.UID, user.LastDynamicTime, time.Duration(conf.DynamicDuration)*time.Minute)
	if err != nil {
		lis.log.WithField("Func", "Listen.Add").Error(err)
		return
	}

	lis.log.WithField("Func", "ListenUpdate.Add").Info(fmt.Sprintf("Listen %d", user.UID))
	go lis.Update.SaveDyanmic(ctx, user.ID, ch)
}

func NewListen(db db.DB, cache cache.Cache, log *logrus.Entry) *ListenUpdate {
	ctx := context.Background()
	comm := listen.New(ctx, listen.NewComment(ctx, conf.GoroutineNum, log), nil, log)
	dyna := listen.New(ctx, listen.NewDynamic(ctx, log), nil, log)

	listen := &ListenUpdate{
		Update: Update{
			comment: comm,
			dynamic: dyna,
			cache:   cache,
			db:      db,
			log:     log,
		},
		Enable: conf.Enable,
		log:    log,
	}

	return listen
}
