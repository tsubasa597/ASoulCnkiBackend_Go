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
}

func (lis ListenUpdate) Started() bool {
	return lis.Enable || atomic.LoadInt32(lis.State) != StateStop
}

func (lis ListenUpdate) Stop() {
	atomic.SwapInt32(lis.State, StateStop)
	lis.Listen.Stop()
	lis.cache.Check.Stop()
	lis.cache.Content.Stop()
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
	go lis.listen(ctx, user.ID, ch)
}

func (lis ListenUpdate) listen(ctx context.Context, userID uint64, ch chan []info.Infoer) {
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
		weight      int64 = 1
		wait, state int32 = 0, 0
		li, ctx           = listen.New(listen.NewDynamic(), &api.API{}, log)
	)

	listen := &ListenUpdate{
		Update: Update{
			State:        &state,
			currentLimit: semaphore.NewWeighted(weight * conf.GoroutineNum),
			weight:       weight,
			Ctx:          ctx,
			log:          log,
			cache:        cache,
			db:           db,
			Wait:         &wait,
			wg:           &sync.WaitGroup{},
		},
		Enable: conf.Enable,
		Listen: li,
		log:    log,
	}

	if listen.Enable {
		atomic.SwapInt32(listen.State, StateRuning)
		listen.getComment(ctx, time.NewTicker(time.Minute*time.Duration(conf.CommentDuration)*2))
	} else {
		atomic.SwapInt32(listen.State, StateStop)
	}

	return listen
}

func (lis ListenUpdate) getComment(ctx context.Context, ticker *time.Ticker) {
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			switch atomic.LoadInt32(lis.State) {
			case StateRuning:
				lis.log.WithField("Func", "listen.Comment").Info("Comment Begin")
				go lis.Comment()
			case StatePause:
				lis.log.WithField("Func", "listen.Comment").Info("Comment Paused")
				continue
			case StateStop:
				lis.log.WithField("Func", "listen.Comment").Info("Comment Stoped")
				return
			}
		}
	}
}
