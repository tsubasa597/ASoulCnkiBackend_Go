package comments

import (
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/api"
)

func Update(user db.User, log *logrus.Entry) {
	var offect int64
	timestamp := user.LastDynamicTime
CommentTag:
	for {
		log.Infoln(db.Update(user))
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
			if info.Time < timestamp {
				break CommentTag
			}
			Add(info.RID, info.CommentType, log)
			offect = info.RID
			user.LastDynamicTime = info.Time
		}
	}
}

func Add(dynamicID int64, commentType uint8, log *logrus.Entry) {
	var (
		wg    = &sync.WaitGroup{}
		ch    = make(chan *db.Comment, 1)
		mutex = &sync.Mutex{}
	)

	for i := 1; true; i++ {
		comment, err := api.GetComments(commentType, dynamicID, conf.DefaultPS, i)
		if err != nil {
			log.Errorln("Func Add api.GetComments Error : ", err)
			continue
		}

		if comment.Code != 0 || len(comment.Data.Replies) == 0 {
			log.Errorln("Func Add Code || Replies Error : ", comment.Message)
			break
		}
		wg.Add(1)
		go Get(dynamicID, comment, wg, ch)
	}

	save := func() {
		for v := range ch {
			if v == nil {
				continue
			}

			mutex.Lock()
			log.Infoln(db.Add(v))
			mutex.Unlock()
		}
	}

	for i := 0; i < 10; i++ {
		go save()
	}

	wg.Wait()
	close(ch)
}

func Get(dynamicID int64, comment *api.Comments, wg *sync.WaitGroup, ch chan<- *db.Comment) {
	for _, comment := range comment.Data.Replies {
		if utf8.RuneCountInString(comment.Content.Message) < conf.DefaultK {
			continue
		}

		ch <- &db.Comment{
			UID:       comment.Mid,
			UName:     comment.Member.Uname,
			Comment:   replacer.Replace(comment.Content.Message),
			DynamicID: dynamicID,
			Time:      comment.Ctime,
		}
	}
	wg.Done()
}

var (
	replacer = strings.NewReplacer("\n", "", " ", "")
)
