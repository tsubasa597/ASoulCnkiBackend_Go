package listen

import (
	"context"
	"math/rand"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/models/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/cache"
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
					Model: entity.Model{
						ID: uint64(comm.Rpid),
					},
					UID:       comm.UID,
					UName:     comm.Name,
					Like:      comm.Like,
					Time:      int64(comm.Time),
					UserID:    userID,
					Content:   comm.Content,
					DynamicID: dynamicID,
				})
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

			err = cache.GetCache().Content.Save()
			if err != nil {
				l.log.WithField("Func", "content.Save").Error(err)
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
					time.Second*time.Duration(rand.Int63n(2)+1))
				if err != nil {
					l.log.WithField("Func", "Listen.Add").Error(err)
					continue
				}

				go l.SaveComment(ct, userID, uint64(dy.DynamicID), ch)
			}
		}
	}
}
