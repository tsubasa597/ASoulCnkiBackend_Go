package entity

import (
	"gorm.io/gorm/clause"
)

// Commentator 评论人表
type Commentator struct {
	Model
	UID       int64  `json:"uid" gorm:"column:uid"`
	Name      string `json:"uname" gorm:"column:uname"`
	LikeNum   uint32 `json:"like" gorm:"column:like_num"`
	Time      int64  `json:"-" gorm:"column:time"`
	DynamicID uint64 `json:"-" gorm:"index:idx_dynamic_id"`
	Rpid      int64  `json:"rpid" gorm:"column:rpid;uniqueIndex"`
	CommentID uint64 `json:"-" gorm:"column:comment_id"`
}

var (
	_ Entity = (*Commentator)(nil)
)

// TableName 表名称
func (Commentator) TableName() string {
	return "commentator"
}

// GetModels 查询时返回的切片
func (Commentator) GetModels() interface{} {
	return &[]Commentator{}
}

// GetClauses 插入时冲突解决方法
func (Commentator) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		Columns:   []clause.Column{{Name: "rpid"}},
		DoUpdates: clause.AssignmentColumns([]string{"like_num", "uname"}),
	}
}
