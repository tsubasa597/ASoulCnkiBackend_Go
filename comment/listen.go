package comment

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
	"github.com/tsubasa597/BILIBILI-HELPER/api"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
	"github.com/tsubasa597/BILIBILI-HELPER/listen"
	"golang.org/x/sync/semaphore"
)

type ListenUpdate struct {
	listen.Listen
	Update
	Duration time.Duration
	log      *logrus.Entry
	dynamic  *listen.Dynamic
	enable   bool
}

func (lis ListenUpdate) Started() bool {
	return lis.enable ||
		(atomic.LoadInt32(lis.Update.commentStarted) == 1 ||
			(atomic.LoadInt32(lis.Update.dynamicStarted) == 1))
}

func (lis ListenUpdate) Add(user entry.User) {
	if !lis.Started() {
		return
	}

	ctx, ch, err := lis.Listen.Add(user.UID, user.LastDynamicTime, lis.dynamic, lis.Duration)
	if err != nil {
		lis.log.WithField("Func", "Listen.Add").Error(err)
		return
	}

	go lis.load(user.UID, user.LastDynamicTime, ctx, ch)

}

func (lis ListenUpdate) load(uid int64, timestamp int32, ctx context.Context, ch <-chan []info.Infoer) {
	for {
		select {
		case <-ctx.Done():
			return
		case infos := <-ch:
			for _, in := range infos {
				dy := in.GetInstance().(info.Dynamic)

				if timestamp >= dy.Time {
					break
				}

				lis.db.Add(&entry.Dynamic{
					UID:       uid,
					DynamicID: dy.DynamicID,
					RID:       dy.RID,
					Type:      dy.CommentType,
					Time:      dy.Time,
					Updated:   false,
				})
			}
		}
	}
}

func NewListen(db db.DB, cache cache.Cacher, log *logrus.Entry) *ListenUpdate {
	var (
		weight    int64 = 1
		dystarted int32 = 0
		costarted int32 = 0
		wait      int32 = 0
		li, ctx         = listen.New(api.API{}, log)
	)

	listen := &ListenUpdate{
		Update: Update{
			currentLimit:   semaphore.NewWeighted(weight * conf.GoroutineNum),
			weight:         weight,
			Ctx:            ctx,
			dynamicStarted: &dystarted,
			commentStarted: &costarted,
			Wait:           &wait,
			wg:             &sync.WaitGroup{},
			cache:          cache,
			db:             db,
			log:            log,
		},
		Duration: time.Duration(time.Minute * time.Duration(conf.Duration)),
		Listen:   *li,
		log:      log,
		dynamic:  listen.NewDynamic(),
		enable:   conf.Satrt,
	}

	return listen
}
