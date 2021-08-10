package comment

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

	atomic.AddInt32(update.dynamicStarted, 1)
	defer atomic.AddInt32(update.dynamicStarted, -1)

	users, err := update.db.Get(&entry.User{})
	if err != nil {
		return
	}

	for _, user := range *users.(*[]entry.User) {
		update.wg.Add(1)
		atomic.AddInt32(update.Wait, 1)
		update.dynamic(user)
	}

	update.wg.Wait()
}

func (update Update) dynamic(user entry.User) {
	update.currentLimit.Acquire(update.Ctx, update.weight)
	var (
		timestamp int32 = user.LastDynamicTime
		offect    int64
		dynamics  = make([]*entry.Dynamic, 0)
	)

DynamicPage:
	for {
		resp, err := api.GetDynamicSrvSpaceHistory(user.UID, offect)
		if err != nil || resp.Data.HasMore != 1 {
			update.log.WithField("Func", "Update api.GetDynamicSrvSpaceHistory").Errorln(err, " ", resp.Message)
			break DynamicPage
		}

		for _, v := range resp.Data.Cards {
			info, err := api.GetOriginCard(v)
			if err != nil {
				update.log.WithField("Func", "Update api.GetOriginCard").Errorln(err, info.CommentType)
				continue
			}

			if info.Time <= timestamp {
				break DynamicPage
			}

			dynamics = append(dynamics, &entry.Dynamic{
				UID:       user.UID,
				DynamicID: info.DynamicID,
				RID:       info.RID,
				Type:      info.CommentType,
				Time:      info.Time,
				Updated:   false,
			})

			user.Name = info.Name
			offect = info.DynamicID
		}
	}

	for i := len(dynamics) - 1; i >= 0; i-- {
		update.db.Add(dynamics[i])

		user.LastDynamicTime = dynamics[i].Time
		update.db.Update(&user, db.Param{
			Field: []string{"dynamic_time", "name"},
		})
	}

	atomic.AddInt32(update.Wait, -1)
	update.currentLimit.Release(update.weight)
	update.wg.Done()
}

func (update *Update) Comment() {
	if atomic.LoadInt32(update.commentStarted) == 1 {
		return
	}

	atomic.AddInt32(update.commentStarted, 1)

	dynamics, err := update.db.Find(&entry.Dynamic{}, db.Param{
		Order: "time asc",
		Query: "is_update = ?",
		Args:  []interface{}{false},
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
	atomic.AddInt32(update.commentStarted, -1)
}

func (update Update) comment(dynamic entry.Dynamic) {
	if dynamic.Updated {
		return
	}

	atomic.AddInt32(update.Wait, 1)
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
			break
		}

		if len(comments.Data.Replies) == 0 {
			update.log.WithField("Func", "UpdateComment Replies").Info(dynamic.Type, dynamic.RID, i)
			break
		}

		comm := make(entry.Comments, 0, len(comments.Data.Replies))
		for _, comment := range comments.Data.Replies {
			s := check.ReplaceStr(comment.Content.Message)

			if utf8.RuneCountInString(s) < conf.DefaultK {
				continue
			}

			comm = append(comm, &entry.Comment{
				UID:       comment.Mid,
				UName:     comment.Member.Uname,
				Comment:   comment.Content.Message,
				CommentID: dynamic.RID,
				Like:      uint32(comment.Like),
				Time:      comment.Ctime,
			})
		}
		update.db.Add(comm)
		update.currentLimit.Release(update.weight)
	}

	dynamic.Updated = true
	update.db.Update(&dynamic, db.Param{
		Field: []string{"is_update"},
	})

	update.currentLimit.Release(update.weight)
	atomic.AddInt32(update.Wait, -1)
	update.wg.Done()
}
