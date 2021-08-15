package entry

import (
	"gorm.io/gorm/clause"
)

type Comment struct {
	Model
	UID       int64  `json:"uid" gorm:"column:uid"`
	UName     string `json:"uname" gorm:"column:uname"`
	Rpid      int64  `json:"rpid" gorm:"column:rpid"`
	Content   string `json:"content" gorm:"column:content"`
	Like      uint32 `json:"like" gorm:"column:like;index:idx_like_time"`
	TotalLike uint32 `json:"total_like" gorm:"column:total_like;index:idx_tl_time"`
	Num       uint32 `json:"num" gorm:"column:num;index:idx_num_time"`
	Time      int64  `json:"comment_time" gorm:"column:time;index:idx_like_time;index:idx_tl_time;index:idx_num_time"`
	DynamicID uint64 `json:"-"`
	UserID    uint64 `json:"-"`
	// Dynamic   *Dynamic `json:"-" gorm:"foreignKey:DynamicID"`
}

var _ Modeler = (*Comment)(nil)

func (Comment) GetModels() interface{} {
	return &[]Comment{}
}

func (Comment) TableName() string {
	return "comment"
}

type Comments []*Comment

var _ Modeler = (*Comments)(nil)

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
