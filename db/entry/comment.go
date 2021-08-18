package entry

import (
	"gorm.io/gorm/clause"
)

type Comment struct {
	Model
	Content   string `json:"content" gorm:"column:content"`
	Like      uint32 `json:"like" gorm:"column:like;index:idx_time"`
	TotalLike uint32 `json:"total_like" gorm:"column:total_like;index:idx_tl"`
	Num       uint32 `json:"num" gorm:"column:num;index:idx_num"`
	Time      int64  `json:"comment_time" gorm:"column:time;index:idx_time"`
	Rpid      int64  `json:"-" gorm:"column:rpid;index:idx_rpid"`
	UserID    uint64 `json:"-"`
}

var (
	_, _ Modeler = (*Comment)(nil), (*Comments)(nil)
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
