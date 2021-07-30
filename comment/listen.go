package comment

import (
	"context"
	"strings"
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
	return lis.enable
}

func (lis Listen) Add(user db.User) {
	if !lis.Started() {
		return
	}

	ctx, ch, err := lis.Listen.Add(user.UID, int64(user.LastDynamicTime), lis.dynamic, lis.Duration)
	if err != nil {
		lis.log.Error("Listen.Add Error: ", err)
		return
	}

	go lis.load(user.UID, ctx, ch)

}

func (lis Listen) load(uid int64, ctx context.Context, ch <-chan []info.Infoer) {
	for {
		select {
		case <-ctx.Done():
			return
		case infos := <-ch:
			for _, in := range infos {
				dy := in.GetInstance().(info.Dynamic)
				lis.chAdd <- db.Dynamic{
					UID:       uid,
					DynamicID: dy.DynamicID,
					RID:       dy.RID,
					Type:      dy.CommentType,
					Time:      dy.Time,
					Updated:   false,
				}
			}
		}
	}
}

func NewListen(log *logrus.Entry) *Listen {
	var weight int64 = 1
	listen := &Listen{
		Update: Update{
			Replacer:     strings.NewReplacer("\n", "", " ", ""),
			CurrentLimit: semaphore.NewWeighted(weight * conf.GoroutineNum),
			Weight:       weight,
			Ctx:          context.Background(),
			Started:      false,
			Wait:         0,
			chAdd:        make(chan db.Modeler, 1),
			chUpdate:     make(chan db.Modeler, 1),
			log:          log,
		},
		Duration: time.Duration(time.Minute * time.Duration(conf.Duration)),
		Listen:   *listen.New(api.API{}, log),
		log:      log,
		dynamic:  listen.NewDynamic(),
		enable:   conf.Satrt,
	}

	// TODO: fix ctx
	go listen.Update.LoadDB(listen.Ctx)

	return listen
}
