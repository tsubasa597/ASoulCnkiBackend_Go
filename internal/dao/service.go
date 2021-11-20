package dao

import (
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/request"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CommentCache 评论缓存，用于查重
type CommentCache struct {
	Rpid    int64
	Content string
}

// GetContent 初始化缓存
func GetContent(batch int, commentCaches []CommentCache, f func(tx *gorm.DB, batch int) error) error {
	return db.Model(&entity.Comment{}).Select("rpid, content").FindInBatches(&commentCaches, batch, f).Error
}

// Rank 作文展数据查询
func Rank(r request.Rank) ([]response.Reply, int64, error) {
	var (
		count   int64
		replies []response.Reply = make([]response.Reply, 0)
	)

	tx := getReply()

	if len(r.Ids) > 0 {
		tx.Where("user.uid in ?", r.Ids)
	}

	if r.Time >= 0 {
		tx.Where("comment.time > ?", r.Time)
	}

	if r.Sort != "" {
		tx.Order(clause.OrderByColumn{
			Column: clause.Column{Name: "comment." + r.Sort},
			Desc:   true,
		})
	}

	tx.Offset((r.Page - 1) * r.Size).Limit(r.Size).Count(&count).Find(&replies)

	return replies, count, tx.Error
}

// Check 查重数据查询
func Check(rpid string) (response.Reply, error) {
	var (
		reply response.Reply
	)

	tx := getReply()
	tx.Where("comment.rpid = ?", rpid).First(&reply)

	return reply, tx.Error
}

// GetTimeInfo 获取时间范围
func GetTimeInfo() (int64, int64) {
	var (
		startTime, endTime int64
	)

	db.Model(&entity.Comment{}).Select("time").Order("time desc").Limit(1).Find(&endTime)
	db.Model(&entity.Comment{}).Select("time").Order("time asc").Limit(1).Find(&startTime)

	return startTime, endTime
}

// getReply 查询条件拼接
func getReply() gorm.DB {
	return *db.Model(&entity.User{}).
		Select("dynamic.type, dynamic.rid as rid, user.uid as uuid, commentator.rpid as rpid, commentator.uid as uid, commentator.time, commentator.uname as name, comment.content, commentator.like_num, commentator.rpid as origin_rpid, comment.num, comment.total_like").
		Joins("inner join comment on comment.dynamic_uid = user.uid").
		Joins("left join commentator on commentator.rpid = comment.rpid").
		Joins("left join dynamic on dynamic.rid = commentator.dynamic_id")
}
