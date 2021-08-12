package entry

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Comment struct {
	Model
	UID       int64  `json:"uid" gorm:"column:uid"`
	UName     string `json:"uname" gorm:"column:uname"`
	Rpid      int64  `json:"rpid" gorm:"column:rpid"`
	Comment   string `json:"comment" gorm:"column:comment"`
	Time      int64  `json:"comment_time" gorm:"column:time"`
	Like      uint32 `json:"like" gorm:"column:like"`
	UserID    uint64
	DynamicID uint64
	Dynamic   *Dynamic `gorm:"foreignKey:DynamicID"`
}

var _ Modeler = (*Comment)(nil)

func (Comment) GetModels() interface{} {
	return &[]Comment{}
}

func (Comment) TableName() string {
	return "comment"
}

func (c *Comment) AfterCreate(tx *gorm.DB) error {
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "comment_text"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"total_like": gorm.Expr("`total_like` + ?", c.Like),
			"num":        gorm.Expr("`num` + 1"),
		}),
	}).Create(&Article{
		Like:        c.Like,
		Time:        c.Time,
		TotalLike:   c.Like,
		Num:         0,
		CommentText: c.Comment,
		CommentID:   c.ID,
		UserID:      c.UserID,
	}).Error; err != nil {
		return err
	}

	return nil
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
