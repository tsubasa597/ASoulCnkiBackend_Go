package update

import (
	"context"
	"sync"
	"sync/atomic"
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

type Update struct {
	Wait           *int32
	Ctx            context.Context
	db             db.DB
	cache          cache.Cacher
	currentLimit   *semaphore.Weighted
	weight         int64
	wg             *sync.WaitGroup
	log            *logrus.Entry
	commentPool    *sync.Pool
	dynamicPool    *sync.Pool
	commentStarted *int32
	dynamicStarted *int32
}

func (update *Update) All() {
	if atomic.LoadInt32(update.commentStarted) == 1 || atomic.LoadInt32(update.dynamicStarted) == 1 {
		return
	}

	update.Dynamic()

	update.Comment()
}

func (update *Update) Dynamic() {
	if atomic.LoadInt32(update.dynamicStarted) == 1 {
		return
	}

	if !atomic.CompareAndSwapInt32(update.dynamicStarted, 0, 1) {
		update.log.WithField("Func", "CompareAndSwapInt32").Error("多协程启动！")
		return
	}
	defer atomic.CompareAndSwapInt32(update.dynamicStarted, 1, 0)

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
		timestamp int32 = user.LastDynamicTime
		offect    int64
		dynamics  = make([]*entry.Dynamic, 0)
	)

DynamicPage:
	for {
		update.currentLimit.Acquire(update.Ctx, update.weight)
		resp, err := api.GetDynamicSrvSpaceHistory(user.UID, offect)
		if err != nil || resp.Data.HasMore != 1 {
			update.log.WithField("Func", "Update api.GetDynamicSrvSpaceHistory").Errorln(err, " ", resp.Message)
			update.currentLimit.Release(update.weight)
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
			})
			user.Name = info.Name
		}
		offect = resp.Data.NextOffset
		update.currentLimit.Release(update.weight)
	}

	for i := len(dynamics) - 1; i >= 0; i-- {
		update.db.Add(dynamics[i])

		user.LastDynamicTime = dynamics[i].Time
		update.db.Update(&user, db.Param{
			Field: []string{"dynamic_time", "name"},
		})
	}
}

func (update *Update) Comment() {
	if atomic.LoadInt32(update.commentStarted) == 1 {
		return
	}

	if !atomic.CompareAndSwapInt32(update.commentStarted, 0, 1) {
		update.log.WithField("Func", "CompareAndSwapInt32").Error("多协程启动！")
		return
	}
	defer atomic.CompareAndSwapInt32(update.commentStarted, 1, 0)

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

	for _, dynamic := range *dynamics.(*[]entry.Dynamic) {
		update.wg.Add(1)
		go update.comment(dynamic)
	}

	update.wg.Wait()
}

func (update Update) comment(dynamic entry.Dynamic) {
	defer update.wg.Done()

	if dynamic.Updated {
		return
	}

	atomic.AddInt32(update.Wait, 1)
	defer atomic.AddInt32(update.Wait, -1)

	for i := 1; true; i++ {
		update.currentLimit.Acquire(update.Ctx, update.weight)
		comments, err := api.GetComments(dynamic.Type, 0, dynamic.RID, conf.DefaultPS, i)
		if err != nil {
			update.log.WithField("Func", "UpdateComment api.GetComments").Errorln(err)
			update.currentLimit.Release(update.weight)
			continue
		}

		if comments.Code != 0 {
			update.log.WithField("Func", "UpdateComment Code").Errorln(comments.Message)
			update.currentLimit.Release(update.weight)
			return
		}

		if len(comments.Data.Replies) == 0 {
			update.log.WithField("Func", "UpdateComment Replies").Info(dynamic.Type, dynamic.RID, i)
			update.currentLimit.Release(update.weight)
			break
		}

		comms := make(entry.Comments, 0, len(comments.Data.Replies))
		for _, comment := range comments.Data.Replies {
			if utf8.RuneCountInString(check.ReplaceStr(comment.Content.Message)) < conf.DefaultK {
				continue
			}

			comm := &entry.Comment{
				UID:       comment.Mid,
				UName:     comment.Member.Uname,
				Comment:   comment.Content.Message,
				Like:      uint32(comment.Like),
				Time:      comment.Ctime,
				Rpid:      comment.Rpid,
				DynamicID: dynamic.ID,
				UserID:    dynamic.UserID,
			}

			if err := update.cache.Increment(*comm, check.HashSet(comment.Content.Message)); err != nil {
				update.log.WithField("Func", "Increment").Error(err)
			}
			comms = append(comms, comm)
		}
		update.currentLimit.Release(update.weight)
		update.db.Add(comms)
	}

	dynamic.Updated = true
	update.db.Update(&dynamic, db.Param{
		Field: []string{"is_update"},
	})

	if err := update.cache.Save(); err != nil {
		update.log.WithField("Func", "cache.Save").Error(err)
	}
}

// func (update Update) rank() {
// 	id, err := update.cache.Get("LastRankID")
// 	if err != nil {
// 		id = "0"
// 	}
// 	comms, err := update.db.Find(&entry.Comment{}, db.Param{
// 		Page:  -1,
// 		Order: "id asc",
// 		Query: "id > ?",
// 		Args:  []interface{}{id},
// 	})
// 	if err != nil {
// 		return
// 	}

// 	for _, comm := range *comms.(*[]entry.Comment) {
// 		update.db.Find(&entry.Article{}, db.Param{
// 			Page:  -1,
// 			Order: "id asc",
// 			Query: "",
// 		})
// 	}
// }
