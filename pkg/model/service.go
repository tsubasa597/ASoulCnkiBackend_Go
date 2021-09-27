package model

import (
	"github.com/tsubasa597/ASoulCnkiBackend/models/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/models/vo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Rank(page, size int, time, order string, uids ...string) (replys []vo.Reply, err error) {
	tx := getReply(page, size)

	if len(uids) > 0 {
		tx.Where("user.uid in ?", uids)
	}

	if time >= "0" {
		tx.Where("comment.time > ?", time)
	}

	if order != "" {
		tx.Order(clause.OrderByColumn{
			Column: clause.Column{Name: "comment." + order},
			Desc:   true,
		})
	}

	tx.Find(&replys)
	return replys, tx.Error
}

func Check(rpid string) (vo.Reply, error) {
	reply := vo.Reply{}

	tx := getReply(2, -1)
	tx.Where("comment.id = ?", rpid).First(&reply)

	return reply, tx.Error
}

type CommentCache struct {
	ID      int64
	Content string
}

func GetContent(rpid string) (commentCache []CommentCache, err error) {
	err = db.Model(&entity.Commentator{}).Select("commentator.id, comment.content").
		Joins("inner join comment comment on commentator.id = comment.id").
		Where("commentator.id > ?", rpid).
		Order("commentator.id asc").
		Find(&commentCache).Error

	return
}

func getReply(page, size int) gorm.DB {
	return *db.Model(&entity.User{}).
		Select("dynamic.type, dynamic.rid as rid, user.uid as uuid, commentator.id as rpid, commentator.uid as uid, commentator.time, commentator.uname as name, comment.content, commentator.like, comment.id as origin_rpid, comment.num, comment.total_like").
		Joins("inner join comment on comment.user_id = user.id").
		Joins("left join dynamic on dynamic.user_id = user.id").
		Joins("left join commentator on commentator.rpid = comment.rpid").
		Offset((page - 1) * size).Limit(size)
}
