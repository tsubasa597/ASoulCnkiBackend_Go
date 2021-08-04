package comment

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/BILIBILI-HELPER/api"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
	"github.com/tsubasa597/BILIBILI-HELPER/listen"
	"golang.org/x/sync/semaphore"
)

type Listen struct {
	listen.Listen
	Update
	Duration time.Duration
	log      *logrus.Entry
	dynamic  *listen.Dynamic
	enable   bool
}

func (lis Listen) Started() bool {
	return lis.enable || (atomic.LoadInt32(lis.Update.Started) == 1)
}

func (lis Listen) Add(user db.User) {
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

func (lis Listen) load(uid int64, timestamp int32, ctx context.Context, ch <-chan []info.Infoer) {
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

				db.Add(&db.Dynamic{
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

func NewListen(log *logrus.Entry) *Listen {
	var (
		weight  int64 = 1
		inx     int32 = 0
		li, ctx       = listen.New(api.API{}, log)
	)

	listen := &Listen{
		Update: Update{
			currentLimit: semaphore.NewWeighted(weight * conf.GoroutineNum),
			weight:       weight,
			Ctx:          ctx,
			Started:      &inx,
			Wait:         &inx,
			wg:           &sync.WaitGroup{},
			log:          log,
		},
		Duration: time.Duration(time.Minute * time.Duration(conf.Duration)),
		Listen:   *li,
		log:      log,
		dynamic:  listen.NewDynamic(),
		enable:   conf.Satrt,
	}

	return listen
}
