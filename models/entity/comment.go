package entity

import (
	"gorm.io/gorm/clause"
)

type Comment struct {
	Model
	Content   string `json:"content" gorm:"column:content"`
	TotalLike uint32 `json:"total_like" gorm:"column:total_like;index:idx_tl_time"`
	Num       uint32 `json:"num" gorm:"column:num;index:idx_num_time"`
	Like      uint32 `json:"like" gorm:"column:like;index:idx_like_time"`
	Time      int64  `json:"comment_time" gorm:"column:time;index:idx_tl_time;index:idx_num_time;index:idx_like_time"`
	UserID    uint64 `json:"-" gorm:"index:idx_user_id"`
}

var (
	_, _ Entity = (*Comment)(nil), (*Comments)(nil)
)

func (Comment) GetModels() interface{} {
	return &[]Comment{}
}

func (Comment) TableName() string {
	return "comment"
}

type Comments []*Comment

func (Comments) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		DoNothing: true,
	}
}

func (Comments) GetModels() interface{} {
	return nil
}

func (Comments) TableName() string {
	return "comment"
}