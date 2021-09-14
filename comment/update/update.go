package update

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/comment/check"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
	"github.com/tsubasa597/BILIBILI-HELPER/listen"
)

type Update struct {
	comment *listen.Listen
	dynamic *listen.Listen
	cache   cache.Cache
	db      db.DB
	log     *logrus.Entry
}

func (u Update) SaveComment(ctx context.Context, userID, dynamicID uint64, ch <-chan []info.Interface) {
	for {
		select {
		case <-ctx.Done():
			return
		case datas := <-ch:
			c := make(entry.Commentators, 0, len(datas))
			for _, data := range datas {
				comm := data.GetInstance().(*info.Comment)

				c = append(c, &entry.Commentator{
					UID:       comm.UID,
					UName:     comm.Name,
					Rpid:      comm.Rpid,
					Like:      comm.Like,
					Time:      int64(comm.Time),
					Content:   comm.Content,
					DynamicID: dynamicID,
					UserID:    userID,
				})

				err := u.cache.Check.Increment(fmt.Sprint(comm.Rpid), check.HashSet(comm.Content))
				if err != nil {
					u.log.WithField("Func", "check.Increment").Error(err)
				}

				err = u.cache.Content.Set(fmt.Sprint(comm.Rpid), comm.Content)
				if err != nil {
					u.log.WithField("Func", "content.Set").Error(err)
				}
			}
			err := u.db.Add(c)
			if err != nil {
				u.log.WithField("Func", "db.Add").Error(err)
			}

			err = u.cache.Check.Save()
			if err != nil {
				u.log.WithField("Func", "check.Save").Error(err)
			}
		}
	}
}

func (u Update) SaveDyanmic(ctx context.Context, userID uint64, ch <-chan []info.Interface) {
	for {
		select {
		case <-ctx.Done():
			return
		case datas := <-ch:
			for _, data := range datas {
				dy := data.GetInstance().(*info.Dynamic)
				u.db.Add(&entry.Dynamic{
					RID:    dy.RID,
					Type:   dy.CommentType,
					Time:   dy.Time,
					Name:   dy.Name,
					UserID: userID,
				})

				ct, ch, err := u.comment.Add(dy.RID, int32(dy.CommentType), time.Second)
				if err != nil {
					u.log.WithField("Func", "db.Add").Error(err)
					continue
				}

				go u.SaveComment(ct, userID, uint64(dy.DynamicID), ch)
			}
		}
	}
}
