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
	*listen.Listen
	Update
	Enable bool
	log    *logrus.Entry
	cancel context.CancelFunc
}

func (lis ListenUpdate) Started() bool {
	return lis.Enable ||
		(atomic.LoadInt32(lis.Update.commentStarted) == 1 ||
			(atomic.LoadInt32(lis.Update.dynamicStarted) == 1))
}

func (lis ListenUpdate) Stop() {
	lis.cancel()
}

func (lis ListenUpdate) Add(user entry.User) {
	if !lis.Started() {
		return
	}

	ctx, ch, err := lis.Listen.Add(user.UID, user.LastDynamicTime, time.Duration(conf.DynamicDuration)*time.Minute)
	if err != nil {
		lis.log.WithField("Func", "Listen.Add").Error(err)
		return
	}

	lis.log.WithField("Func", "ListenUpdate.Add").Info(fmt.Sprintf("Listen %d", user.UID))
	go lis.listen(user.ID, ch, ctx)
}

func (lis ListenUpdate) listen(userID uint64, ch chan []info.Infoer, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case infos := <-ch:
			for i := len(infos) - 1; i >= 0; i-- {
				dy := infos[i].GetInstance().(*info.Dynamic)
				if err := lis.db.Add(&entry.Dynamic{
					RID:     dy.RID,
					Type:    dy.CommentType,
					Time:    dy.Time,
					Updated: false,
					Name:    dy.Name,
					UserID:  userID,
				}); err != nil {
					lis.log.WithField("Func", "db.Add").Error(err)
				}
			}
		}
	}
}

func NewListen(db db.DB, cache cache.Cache, log *logrus.Entry) *ListenUpdate {
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
		Listen: li,
		log:    log,
		Enable: conf.Enable,
	}

	if conf.Enable {
		go func(ticker <-chan time.Time, ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker:
					listen.log.WithField("Func", "listen.Comment").Info("Comment Begin")
					go listen.Comment()
				}
			}
		}(time.NewTicker(time.Duration(conf.CommentDuration)*time.Minute).C, ctx)
	}

	return listen
}
