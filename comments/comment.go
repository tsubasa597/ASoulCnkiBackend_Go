package comments

import (
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/api"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
)

// TODO: 导入缓存
func Update(user db.User, log *logrus.Entry) {
	if started {
		return
	}

	started = true
	var (
		offect    int64
		timestamp = user.LastDynamicTime
		q         = make(queue, 0)
		ch        = make(chan db.Comments, 1)
	)

	wg.Add(1)
	go func() {
	CommentTag:
		for {
			resp, err := api.GetDynamicSrvSpaceHistory(user.UID, offect)
			if err != nil {
				log.Errorln("Func Update api.GetDynamicSrvSpaceHistory Error : ", err)
				return
			}

			for _, v := range resp.Data.Cards {
				info, err := api.GetOriginCard(v)
				if err != nil {
					log.Errorln("Func Update api.GetOriginCard Error : ", err)
					return
				}

				if info.Time <= timestamp {
					break CommentTag
				}

				q.push(&info)
				wait++
				offect = info.RID
			}
		}

		for info := q.pop(); info != nil; info = q.pop() {
			wg.Add(1)
			go Add(info.RID, info.CommentType, ch, log)

			user.LastDynamicTime = info.Time
			log.Infoln("Update User Error: ", db.Update(user))
		}
		wg.Done()
	}()

	go func() {
		for v := range ch {
			log.Info("Add Comment Error : ", db.Add(v))
		}
	}()

	started = false
	wg.Wait()
	close(ch)
}

func Add(commentID int64, commentType uint8, ch chan<- db.Comments, log *logrus.Entry) {
	for i := 1; true; i++ {
		comments, err := api.GetComments(commentType, commentID, conf.DefaultPS, i)
		if err != nil {
			log.Errorln("Func Add api.GetComments Error : ", err)
			continue
		}

		if comments.Code != 0 || len(comments.Data.Replies) == 0 {
			log.Errorln("Func Add Code || Replies Error : ", comments.Message)
			break
		}

		comm := make(db.Comments, 0, len(comments.Data.Replies))
		for _, comment := range comments.Data.Replies {
			if utf8.RuneCountInString(comment.Content.Message) < conf.DefaultK {
				continue
			}

			s := replacer.Replace(comment.Content.Message)
			var emote string
			for k, v := range comment.Content.Emote {
				s = strings.Replace(s, k, string(v.Id), -1)
				emote = k + ","

				err = db.Find(&db.Emote{
					EmoteID:   v.Id,
					EmoteText: k,
				})
				if err != nil && err.Error() == db.ErrNotFound {
					log.Info("Add Emote Error : ", db.Add(&db.Emote{
						EmoteID:   v.Id,
						EmoteText: k,
					}))
				}

			}

			comm = append(comm, &db.Comment{
				UID:       comment.Mid,
				UName:     comment.Member.Uname,
				Comment:   s,
				CommentID: commentID,
				Time:      comment.Ctime,
				Emote:     emote,
			})
		}
		ch <- comm
	}

	wg.Done()
	wait--
}

func Status() (bool, int) {
	return started, wait
}

var (
	replacer = strings.NewReplacer("\n", "", " ", "")
	started  bool
	wait     int
	wg       = &sync.WaitGroup{}
)

type queue []*info.Dynamic

func (q *queue) pop() *info.Dynamic {
	n := len(*q)
	if n != 0 {
		d := (*q)[n-1]
		*q = (*q)[:n-1]
		return d
	} else {
		return nil
	}
}

func (q *queue) push(dynamic *info.Dynamic) {
	*q = append(*q, dynamic)
}
