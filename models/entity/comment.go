package entity

import (
	"gorm.io/gorm/clause"
)

// Comment 评论表
type Comment struct {
	Model
	Content   string `json:"content" gorm:"column:content"`
	TotalLike uint32 `json:"total_like" gorm:"column:total_like;index:idx_tl_time"`
	Num       uint32 `json:"num" gorm:"column:num;index:idx_num_time"`
	Like      uint32 `json:"like" gorm:"column:like;index:idx_like_time"`
	Time      int64  `json:"comment_time" gorm:"column:time;index:idx_tl_time;index:idx_num_time;index:idx_like_time"`
	UserID    uint64 `json:"-" gorm:"index:idx_user_id"`
	Rpid      int64  `json:"rpid" gorm:"column:rpid;index:idx_rpid"`
}

var (
	_, _ Entity = (*Comment)(nil), (*Comments)(nil)
)

// GetModels 查询时返回的切片
func (Comment) GetModels() interface{} {
	return &[]Comment{}
}

// TableName 表名称
func (Comment) TableName() string {
	return "comment"
}

// Comments 批量插入使用
type Comments []*Comment

// GetClauses 插入时冲突解决方法
func (Comments) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		DoNothing: true,
	}
}

// GetModels 查询时返回的切片
func (Comments) GetModels() interface{} {
	return nil
}

// TableName 表名称
func (Comments) TableName() string {
	return "comment"
}
