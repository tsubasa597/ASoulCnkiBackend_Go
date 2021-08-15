package update

import (
	"context"
	"fmt"
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
	Enable   bool
	Duration time.Duration
	log      *logrus.Entry
}

func (lis ListenUpdate) Started() bool {
	return lis.Enable ||
		(atomic.LoadInt32(lis.Update.commentStarted) == 1 ||
			(atomic.LoadInt32(lis.Update.dynamicStarted) == 1))
}

func (lis ListenUpdate) Add(user entry.User) {
	if !lis.Started() {
		return
	}

	ctx, ch, err := lis.Listen.Add(user.UID, user.LastDynamicTime, lis.Duration)
	if err != nil {
		lis.log.WithField("Func", "Listen.Add").Error(err)
		return
	}

	lis.log.WithField("Func", "ListenUpdate.Add").Info(fmt.Sprintf("Listen %d", user.UID))
	go lis.load(&user, user.LastDynamicTime, ctx, ch)
}

func (lis ListenUpdate) load(user *entry.User, timestamp int32, ctx context.Context, ch <-chan []info.Infoer) {
	for {
		select {
		case <-ctx.Done():
			return
		case infos := <-ch:
			for _, in := range infos {
				dy := in.GetInstance().(*info.Dynamic)

				if timestamp >= dy.Time {
					break
				}

				dynaimc := &entry.Dynamic{
					UserID:  user.ID,
					RID:     dy.RID,
					Type:    dy.CommentType,
					Time:    dy.Time,
					Updated: false,
				}

				lis.db.Add(dynaimc)

				lis.wg.Add(1)
				go lis.comment(*dynaimc)
			}
		}
	}
}

func NewListen(db db.DB, cache cache.Cacher, log *logrus.Entry) *ListenUpdate {
	var (
		weight                     int64 = 1
		dystarted, costarted, wait int32 = 0, 0, 0
		li, ctx                          = listen.New(listen.NewDynamic(), &api.API{}, log)
	)

	listen := &ListenUpdate{
		Update: Update{
			currentLimit:   semaphore.NewWeighted(weight * conf.GoroutineNum),
			weight:         weight,
			Ctx:            ctx,
			log:            log,
			cache:          cache,
			db:             db,
			dynamicStarted: &dystarted,
			commentStarted: &costarted,
			Wait:           &wait,
			wg:             &sync.WaitGroup{},
		},
		Duration: time.Duration(time.Minute * time.Duration(conf.Duration)),
		Listen:   *li,
		log:      log,
		Enable:   conf.Satrt,
	}

	return listen
}
