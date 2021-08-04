package comment

import (
	"context"
	"sync"
	"sync/atomic"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"golang.org/x/sync/semaphore"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/api"
)

type Update struct {
	CurrentLimit *semaphore.Weighted
	Weight       int64
	Ctx          context.Context
	Started      bool
	Wait         *int32
	log          *logrus.Entry
}

func (update *Update) Do(users []db.User) {
	if update.Started {
		return
	}

	update.Started = true

	var (
		wg = &sync.WaitGroup{}
	)

	for _, user := range users {
		wg.Add(1)
		atomic.AddInt32(update.Wait, 1)
		go update.UpdateDynamic(user, wg)
	}
	wg.Wait()

	for _, dynamic := range *(db.Dynamic{}.Find([]interface{}{"is_update <> true"})).(*[]db.Dynamic) {
		wg.Add(1)
		atomic.AddInt32(update.Wait, 1)
		go update.UpdateComment(dynamic, wg)
	}
	wg.Wait()

	update.Started = false
}

func (update Update) UpdateDynamic(user db.User, wg *sync.WaitGroup) {
	update.CurrentLimit.Acquire(update.Ctx, update.Weight)
	var (
		timestamp int32 = user.LastDynamicTime
		offect    int64
		dynamics  = make([]*db.Dynamic, 0)
	)

	for {
		resp, err := api.GetDynamicSrvSpaceHistory(user.UID, offect)
		if err != nil || resp.Data.HasMore != 1 {
			update.log.WithField("Func", "Update api.GetDynamicSrvSpaceHistory").Errorln(err, " ", resp.Message)
			break
		}

		for _, v := range resp.Data.Cards {
			info, err := api.GetOriginCard(v)
			if err != nil {
				update.log.WithField("Func", "Update api.GetOriginCard").Errorln(err, info.CommentType)
				continue
			}

			if info.Time <= timestamp {
				break
			}

			dynamics = append(dynamics, &db.Dynamic{
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
		db.Add(dynamics[i])

		user.LastDynamicTime = dynamics[i].Time
		db.Update(&user)
	}

	atomic.AddInt32(update.Wait, -1)
	update.CurrentLimit.Release(update.Weight)
	wg.Done()
}

func (update Update) UpdateComment(dynamic db.Dynamic, wg *sync.WaitGroup) {
	update.CurrentLimit.Acquire(update.Ctx, update.Weight)

	for i := 1; true; i++ {
		comments, err := api.GetComments(dynamic.Type, 0, dynamic.RID, conf.DefaultPS, i)
		if err != nil {
			update.log.WithField("Func", "UpdateComment api.GetComments").Errorln(err)
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

		comm := make(db.Comments, 0, len(comments.Data.Replies))
		for _, comment := range comments.Data.Replies {
			s := ReplaceStr(comment.Content.Message)

			if utf8.RuneCountInString(s) < conf.DefaultK {
				continue
			}

			comm = append(comm, &db.Comment{
				UID:       comment.Mid,
				UName:     comment.Member.Uname,
				Comment:   s,
				CommentID: dynamic.RID,
				Like:      uint32(comment.Like),
				Time:      comment.Ctime,
			})
		}
		db.Add(comm)
	}

	dynamic.Updated = true
	db.Update(&dynamic)

	atomic.AddInt32(update.Wait, -1)
	update.CurrentLimit.Release(update.Weight)
	wg.Done()
}
