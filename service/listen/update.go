package listen

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/models/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/check"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/model"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
)

// SaveComment 更新评论区评论
func (l Listen) SaveComment(ctx context.Context, userID, dynamicID uint64, ch <-chan []info.Interface) {
	for {
		select {
		case <-ctx.Done():
			return
		case datas := <-ch:
			c := make(entity.Commentators, 0, len(datas))
			for _, data := range datas {
				comm := data.GetInstance().(*info.Comment)

				c = append(c, &entity.Commentator{
					UID:       comm.UID,
					UName:     comm.Name,
					Rpid:      comm.Rpid,
					Like:      comm.Like,
					Time:      int64(comm.Time),
					Content:   comm.Content,
					DynamicID: dynamicID,
					UserID:    userID,
				})

				err := cache.GetCache().Check.Increment(fmt.Sprint(comm.Rpid), check.HashSet(comm.Content))
				if err != nil {
					l.log.WithField("Func", "check.Increment").Error(err)
				}

				err = cache.GetCache().Content.Set(fmt.Sprint(comm.Rpid), comm.Content)
				if err != nil {
					l.log.WithField("Func", "content.Set").Error(err)
				}
			}

			if len(c) == 0 {
				continue
			}

			err := model.Add(c)
			if err != nil {
				l.log.WithField("Func", "db.Add").Error(err)
			}

			err = cache.GetCache().Check.Save()
			if err != nil {
				l.log.WithField("Func", "check.Save").Error(err)
			}
		}
	}
}

// SaveDyanmic 更新动态并更新评论区评论
func (l Listen) SaveDyanmic(ctx context.Context, userID uint64, ch <-chan []info.Interface) {
	for {
		select {
		case <-ctx.Done():
			return
		case datas := <-ch:
			for _, data := range datas {
				dy := data.GetInstance().(*info.Dynamic)
				model.Add(&entity.Dynamic{
					RID:    dy.RID,
					Type:   dy.CommentType,
					Time:   dy.Time,
					Name:   dy.Name,
					UserID: userID,
				})

				ct, ch, err := l.comment.Add(dy.RID, int32(dy.CommentType),
					time.Second*time.Duration(rand.Intn(4)+1))
				if err != nil {
					l.log.WithField("Func", "db.Add").Error(err)
					continue
				}

				go l.SaveComment(ct, userID, uint64(dy.DynamicID), ch)
			}
		}
	}
}
