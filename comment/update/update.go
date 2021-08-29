package update

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/comment/check"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
	"golang.org/x/sync/semaphore"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/api"
)

const (
	StateStop int32 = iota + 1
	StateRuning
	StatePause
)

type Update struct {
	Wait         *int32
	State        *int32
	Ctx          context.Context
	db           db.DB
	cache        cache.Cache
	currentLimit *semaphore.Weighted
	weight       int64
	wg           *sync.WaitGroup
	log          *logrus.Entry
}

func (update *Update) All() {
	if atomic.LoadInt32(update.State) != StateStop {
		return
	}

	update.Dynamic()

	update.Comment()
}

func (update *Update) Dynamic() {
	if atomic.LoadInt32(update.State) != StateStop {
		return
	}

	if !atomic.CompareAndSwapInt32(update.State, StateStop, StateRuning) {
		update.log.WithField("Func", "CompareAndSwapInt32").Error("多协程启动！")
		return
	}
	defer atomic.CompareAndSwapInt32(update.State, StateRuning, StateStop)

	users, err := update.db.Get(&entry.User{})
	if err != nil {
		return
	}

	for _, user := range *users.(*[]entry.User) {
		update.wg.Add(1)
		go update.dynamic(user)
	}

	update.wg.Wait()
}

func (update Update) dynamic(user entry.User) {
	defer update.wg.Done()

	atomic.AddInt32(update.Wait, 1)
	defer atomic.AddInt32(update.Wait, -1)

	var (
		dynamics  entry.Dynamics = make(entry.Dynamics, 0)
		timestamp int32          = user.LastDynamicTime
		offect    int64          = 0
	)

DynamicPage:
	for {
		if atomic.LoadInt32(update.State) != StateRuning {
			update.log.WithField("Func", "LoadInt32").Infof("State changed: %d", atomic.LoadInt32(update.State))
			return
		}

		update.currentLimit.Acquire(update.Ctx, update.weight)
		resp, err := api.GetDynamicSrvSpaceHistory(user.UID, offect)
		if err != nil || resp.Data.HasMore != 1 {
			update.log.WithField("Func", "Update api.GetDynamicSrvSpaceHistory").Errorln(err, " ", resp.Message)
			update.currentLimit.Release(update.weight)
			break DynamicPage
		}

		if resp.Code != 0 {
			update.log.WithField("Func", "UpdateComment Code").Errorln(resp.Message)
			update.Pause()
			break DynamicPage
		}

		for _, v := range resp.Data.Cards {
			info, err := api.GetOriginCard(v)
			if err != nil {
				update.log.WithField("Func", "Update api.GetOriginCard").Errorln(err)
				continue
			}

			if info.Time <= timestamp {
				update.currentLimit.Release(update.weight)
				break DynamicPage
			}

			dynamics = append(dynamics, &entry.Dynamic{
				UserID:  user.ID,
				RID:     info.RID,
				Type:    info.CommentType,
				Time:    info.Time,
				Updated: false,
				Name:    info.Name,
			})
		}
		offect = resp.Data.NextOffset
		update.currentLimit.Release(update.weight)
	}

	sort.Sort(dynamics)
	if err := update.db.Add(dynamics); err != nil {
		update.log.WithField("Func", "db.Add").Error(err)
	}
}

func (update *Update) Comment() {
	if atomic.LoadInt32(update.State) != StateStop {
		return
	}

	if !atomic.CompareAndSwapInt32(update.State, StateStop, StateRuning) {
		update.log.WithField("Func", "CompareAndSwapInt32").Error("多协程启动！")
		return
	}
	defer atomic.CompareAndSwapInt32(update.State, StateRuning, StateStop)

	dynamics, err := update.db.Find(&entry.Dynamic{}, db.Param{
		Order: "time asc",
		Query: "is_update = ?",
		Args:  []interface{}{false},
		Page:  -1,
	})
	if err != nil {
		update.log.WithField("Func", "DB.Find").Error(err)
		return
	}

	atomic.AddInt32(update.Wait, int32(len(*dynamics.(*[]entry.Dynamic))))
	for _, dynamic := range *dynamics.(*[]entry.Dynamic) {
		update.wg.Add(1)
		update.comment(dynamic)
	}

	update.wg.Wait()
}

func (update Update) comment(dynamic entry.Dynamic) {
	defer update.wg.Done()
	defer atomic.AddInt32(update.Wait, -1)

	if dynamic.Updated {
		return
	}

	for i := 1; true; i++ {
		update.currentLimit.Acquire(update.Ctx, update.weight)

		if atomic.LoadInt32(update.State) != StateRuning {
			update.log.WithField("Func", "LoadInt32").Infof("State changed: %d", atomic.LoadInt32(update.State))
			return
		}

		comments, err := api.GetComments(dynamic.Type, 0, dynamic.RID, conf.DefaultPS, i)
		if err != nil {
			update.log.WithField("Func", "UpdateComment api.GetComments").Errorln(err)
			update.currentLimit.Release(update.weight)
			continue
		}

		if comments.Code != 0 {
			update.log.WithField("Func", "UpdateComment Code").Errorln(comments.Message)
			update.Pause()
			update.currentLimit.Release(update.weight)
			return
		}

		if len(comments.Data.Replies) == 0 {
			update.log.WithField("Func", "UpdateComment Replies").Info(dynamic.Type, dynamic.RID, i)
			update.currentLimit.Release(update.weight)
			break
		}
		update.currentLimit.Release(update.weight)

		comms := make(entry.Commentators, 0, len(comments.Data.Replies))
		for _, comment := range comments.Data.Replies {
			if utf8.RuneCountInString(check.ReplaceStr(comment.Content.Message)) < conf.DefaultK {
				continue
			}

			comm := &entry.Commentator{
				UID:       comment.Mid,
				UName:     comment.Member.Uname,
				Content:   comment.Content.Message,
				Like:      uint32(comment.Like),
				Time:      comment.Ctime,
				Rpid:      comment.Rpid,
				DynamicID: dynamic.ID,
				UserID:    dynamic.UserID,
			}

			if err := update.cache.Check.Increment(fmt.Sprint(comment.Rpid),
				check.HashSet(comment.Content.Message)); err != nil {
				update.log.WithField("Func", "cache.Content.Increment").Error(err)
			}
			update.cache.Content.Set(fmt.Sprint(comment.Rpid), comment.Content.Message)

			comms = append(comms, comm)
		}
		if err := update.db.Add(comms); err != nil {
			update.log.WithField("Func", "db.Add").Error(err)
		}
	}

	dynamic.Updated = true
	update.db.Update(&dynamic, db.Param{
		Field: []string{"is_update"},
	})

	if err := update.cache.Check.Save(); err != nil {
		update.log.WithField("Func", "cache.Save").Error(err)
	}
}

func (update Update) Pause() {
	update.log.WithField("Func", "update.Pause").Info("接口限流，暂停爬取")

	if atomic.CompareAndSwapInt32(update.State, StateRuning, StatePause) {
		go func(ctx context.Context, state *int32, ticker *time.Ticker) {
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					atomic.CompareAndSwapInt32(state, StatePause, StateRuning)
					return
				}
			}
		}(update.Ctx, update.State, time.NewTicker(time.Hour))
	}
}
