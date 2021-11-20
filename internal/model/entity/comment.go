package entity

import (
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

// Comment 评论表
type Comment struct {
	Model
	Content    string `json:"content" gorm:"column:content"`
	TotalLike  uint32 `json:"total_like" gorm:"column:total_like;index:idx_tl_time"`
	Num        uint32 `json:"num" gorm:"column:num;index:idx_num_time"`
	LikeNum    uint32 `json:"like" gorm:"column:like_num;index:idx_like_time"`
	Time       int64  `json:"comment_time" gorm:"column:time;index:idx_tl_time;index:idx_num_time;index:idx_like_time"`
	DynamicUID int64  `json:"-" gorm:"index:idx_dy_uid"`
	DynamicID  uint64 `json:"-" gorm:"column:dynamic_id"`
	Rpid       int64  `json:"rpid" gorm:"column:rpid"`
	UID        int64  `json:"-" gorm:"-"`
	Name       string `json:"-" gorm:"-"`
}

var (
	_, _            Entity                          = (*Comment)(nil), (*Comments)(nil)
	_               callbacks.BeforeCreateInterface = (*Comment)(nil)
	commentatorPool sync.Pool                       = sync.Pool{
		New: func() interface{} {
			return &Commentator{}
		},
	}
	commentPool sync.Pool = sync.Pool{
		New: func() interface{} {
			return &Comment{}
		},
	}
)

// GetModels 查询时返回的切片
func (Comment) GetModels() interface{} {
	return &[]Comment{}
}

// TableName 表名称
func (Comment) TableName() string {
	return "comment"
}

// GetClauses 插入时冲突解决方法
func (Comment) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"time", "rpid", "like_num", "dynamic_uid", "dynamic_id", "total_like", "num",
		}),
	}
}

// BeforeCreate 插入之前进行部分数据更改
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	var (
		comment     = commentPool.Get().(*Comment)
		commentator = commentatorPool.Get().(*Commentator)
		likeNum     uint32
	)
	defer func() {
		commentatorPool.Put(commentator)
		commentPool.Put(comment)
	}()

	tx.Model(commentator).Select("like_num").Where("rpid", c.Rpid).Find(&likeNum)

	commentator.DynamicID = c.DynamicID
	commentator.LikeNum = c.LikeNum
	commentator.Rpid = c.Rpid
	commentator.Time = c.Time
	commentator.UID = c.UID
	commentator.Name = c.Name
	commentator.CommentID = c.ID

	tx.Model(commentator).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "rpid"}},
		DoUpdates: clause.AssignmentColumns([]string{"like_num"}),
	}).Select("dynamic_id", "like_num", "rpid", "time", "uid", "uname", "comment_id").Create(commentator)

	if tx.Model(comment).Select("id", "time", "rpid", "like_num", "dynamic_uid", "dynamic_id", "num", "total_like").
		Where("content = ?", c.Content).Find(comment).RowsAffected != 0 {

		c.ID = comment.ID
		c.TotalLike = comment.TotalLike + c.LikeNum - likeNum
		c.Num = comment.Num + 1
		if c.Time > comment.Time {
			c.Time = comment.Time
			c.Rpid = comment.Rpid
			c.LikeNum = comment.LikeNum
			c.DynamicUID = comment.DynamicUID
			c.DynamicID = comment.DynamicID
		}
	} else {
		c.TotalLike = c.LikeNum
	}

	return nil
}

type Comments []Comment

// GetModels 查询时返回的切片
func (Comments) GetModels() interface{} {
	return &[]Comment{}
}

// TableName 表名称
func (Comments) TableName() string {
	return "comment"
}

// GetClauses 插入时冲突解决方法
func (Comments) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"time", "rpid", "like_num", "dynamic_uid", "dynamic_id", "total_like", "num",
		}),
	}
}
