package comment

import (
	"context"
	"strings"
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
	Replacer     *strings.Replacer
	CurrentLimit *semaphore.Weighted
	Weight       int64
	Ctx          context.Context
	Started      bool
	Wait         int32
	chAdd        chan db.Modeler
	chUpdate     chan db.Modeler
	log          *logrus.Entry
}

func (update *Update) Do(users []db.User) {
	if update.Started {
		return
	}

	update.Started = true

	var (
		canCtx, cancle = context.WithCancel(update.Ctx)
		wg             = &sync.WaitGroup{}
	)

	go update.LoadDB(canCtx)

	for _, user := range users {
		wg.Add(1)
		go update.UpdateDynamic(user, update.chAdd, update.chUpdate, wg)
	}
	wg.Wait()

	for _, dynamic := range *(db.Dynamic{}.Find([]interface{}{"is_update <> 1"})).(*[]db.Dynamic) {
		wg.Add(1)
		go update.UpdateComment(dynamic, update.chAdd, update.chUpdate, wg)
	}
	wg.Wait()

	LoadCache()
	update.Started = false
	cancle()
}

func (update Update) UpdateDynamic(user db.User, chAdd, chUpdate chan<- db.Modeler, wg *sync.WaitGroup) {
	atomic.AddInt32(&update.Wait, 1)
	update.CurrentLimit.Acquire(update.Ctx, update.Weight)
	var (
		timestamp int32 = user.LastDynamicTime
		offect    int64
	)

	for {
		resp, err := api.GetDynamicSrvSpaceHistory(user.UID, offect)
		if err != nil || resp.Data.HasMore != 1 {
			update.log.Errorln("Func Update api.GetDynamicSrvSpaceHistory Error : ", err, " ", resp.Message)

			atomic.AddInt32(&update.Wait, -1)
			update.CurrentLimit.Release(update.Weight)
			wg.Done()

			return
		}

		for _, v := range resp.Data.Cards {
			info, err := api.GetOriginCard(v)
			if err != nil {
				update.log.Errorln("Func Update api.GetOriginCard Error : ", err, info.CommentType)
				continue
			}

			if info.Time <= timestamp {
				atomic.AddInt32(&update.Wait, -1)
				update.CurrentLimit.Release(update.Weight)
				wg.Done()

				return
			}

			chAdd <- &db.Dynamic{
				UID:       user.UID,
				DynamicID: info.DynamicID,
				RID:       info.RID,
				Type:      info.CommentType,
				Time:      info.Time,
				Updated:   false,
			}

			user.LastDynamicTime = info.Time
			user.Name = info.Name
			chUpdate <- &user
			offect = info.DynamicID
		}
	}
}

func (update Update) UpdateComment(dynamic db.Dynamic, chAdd, chUpdate chan<- db.Modeler, wg *sync.WaitGroup) {
	atomic.AddInt32(&update.Wait, 1)
	update.CurrentLimit.Acquire(update.Ctx, update.Weight)

	for i := 1; true; i++ {
		comments, err := api.GetComments(dynamic.Type, 0, dynamic.RID, conf.DefaultPS, i)
		if err != nil {
			update.log.Errorln("Func Add api.GetComments Error : ", err)
			continue
		}

		if comments.Code != 0 || len(comments.Data.Replies) == 0 {
			update.log.Errorln("Func Add Code || Replies Error : ", comments.Message, dynamic.Type, dynamic.RID, i)
			break
		}

		comm := make(db.Comments, 0, len(comments.Data.Replies))
		for _, comment := range comments.Data.Replies {
			s := update.Replacer.Replace(comment.Content.Message)

			for k, v := range comment.Content.Emote {
				s = strings.Replace(s, k, string(v.Id), -1)

				chAdd <- &db.Emote{
					EmoteID:   v.Id,
					EmoteText: k,
				}
			}

			if utf8.RuneCountInString(s) < conf.DefaultK {
				continue
			}

			comm = append(comm, &db.Comment{
				UID:       comment.Mid,
				UName:     comment.Member.Uname,
				Comment:   s,
				CommentID: dynamic.RID,

				Time: comment.Ctime,
			})
		}
		chAdd <- comm
	}

	dynamic.Updated = true
	chUpdate <- &dynamic

	atomic.AddInt32(&update.Wait, -1)
	update.CurrentLimit.Release(update.Weight)
	wg.Done()
}

func (update Update) LoadDB(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-update.chAdd:
			if err := db.Add(data); err != nil {
				update.log.Error("Add Error: ", err)
			}

		case data := <-update.chUpdate:
			if err := db.Update(data); err != nil {
				update.log.Error("Update Error: ", err)
			}
		}
	}
}
