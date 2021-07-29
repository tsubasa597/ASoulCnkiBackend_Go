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

var (
	replacer           = strings.NewReplacer("\n", "", " ", "")
	currentLimit       = semaphore.NewWeighted(weight * 10)
	ctx                = context.Background()
	weight       int64 = 1
	started      bool
	wait         int32
)

func Update(log *logrus.Entry) {
	if started {
		return
	}

	started = true

	var (
		canCtx, cancle  = context.WithCancel(ctx)
		chAdd, chUpdate = make(chan db.Modeler), make(chan db.Modeler)
		wg              = &sync.WaitGroup{}
	)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-chAdd:
				log.Info("Add Error: ", db.Add(data))
			case data := <-chUpdate:
				log.Info("Update Error: ", db.Update(data))
			}

		}
	}(canCtx)

	users := *db.Get(&db.User{}).(*[]db.User)

	for _, user := range users {
		wg.Add(1)
		go UpdateDynamic(user, chAdd, chUpdate, wg, log)
	}
	wg.Wait()

	for _, dynamic := range *(db.Dynamic{}.Find([]interface{}{"is_update <> 1"})).(*[]db.Dynamic) {
		wg.Add(1)
		go UpdateComment(dynamic, chAdd, chUpdate, wg, log)
	}
	wg.Wait()

	LoadCache()
	started = false
	cancle()
	defer close(chAdd)
	defer close(chUpdate)
}

func UpdateDynamic(user db.User, chAdd chan<- db.Modeler, chUpdate chan<- db.Modeler, wg *sync.WaitGroup, log *logrus.Entry) {
	atomic.AddInt32(&wait, 1)
	currentLimit.Acquire(ctx, weight)
	var (
		timestamp int32 = user.LastDynamicTime
		offect    int64
	)

	for {
		resp, err := api.GetDynamicSrvSpaceHistory(user.UID, offect)
		if err != nil || resp.Data.HasMore != 1 {
			log.Errorln("Func Update api.GetDynamicSrvSpaceHistory Error : ", err, " ", resp.Message)

			atomic.AddInt32(&wait, -1)
			currentLimit.Release(weight)
			wg.Done()

			return
		}

		for _, v := range resp.Data.Cards {
			info, err := api.GetOriginCard(v)
			if err != nil {
				log.Errorln("Func Update api.GetOriginCard Error : ", err, info.CommentType)
				continue
			}

			if info.Time <= timestamp {
				atomic.AddInt32(&wait, -1)
				currentLimit.Release(weight)
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

func UpdateComment(dynamic db.Dynamic, chAdd chan<- db.Modeler, chUpdate chan<- db.Modeler, wg *sync.WaitGroup, log *logrus.Entry) {
	atomic.AddInt32(&wait, 1)
	currentLimit.Acquire(ctx, weight)

	for i := 1; true; i++ {
		comments, err := api.GetComments(dynamic.Type, 0, dynamic.RID, conf.DefaultPS, i)
		if err != nil {
			log.Errorln("Func Add api.GetComments Error : ", err)
			continue
		}

		if comments.Code != 0 || len(comments.Data.Replies) == 0 {
			log.Errorln("Func Add Code || Replies Error : ", comments.Message, len(comments.Data.Replies), dynamic.Type, dynamic.RID, i)
			break
		}

		comm := make(db.Comments, 0, len(comments.Data.Replies))
		for _, comment := range comments.Data.Replies {
			s := replacer.Replace(comment.Content.Message)

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

	atomic.AddInt32(&wait, -1)
	currentLimit.Release(weight)
	wg.Done()
}

func Status() (bool, int32) {
	return started, wait
}
