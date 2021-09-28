package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/pkg/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/check"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

// Commentator 评论人表
type Commentator struct {
	Model
	UID       int64  `json:"uid" gorm:"column:uid"`
	UName     string `json:"uname" gorm:"column:uname"`
	Like      uint32 `json:"like" gorm:"column:like"`
	Time      int64  `json:"-" gorm:"column:time"`
	Content   string `json:"-" gorm:"-"`
	UserID    uint64 `json:"-" grom:"-"`
	DynamicID uint64 `json:"-" gorm:"index:idx_dynamic_id"`
	CommentID uint64 `json:"-"`
	Rpid      int64  `json:"rpid" gorm:"column:rpid;index:idx_rpid"`
}

var (
	_           Entity                          = (*Commentator)(nil)
	_           callbacks.BeforeCreateInterface = (*Commentator)(nil)
	commentPool sync.Pool                       = sync.Pool{
		New: func() interface{} {
			return &Comment{}
		},
	}
)

// TableName 表名称
func (Commentator) TableName() string {
	return "commentator"
}

// GetModels 查询时返回的切片
func (Commentator) GetModels() interface{} {
	return &[]Commentator{}
}

// BeforeCreate 在插入之前更新评论表中的数据
func (c *Commentator) BeforeCreate(tx *gorm.DB) error {
	comm := commentPool.Get().(*Comment)
	defer func() {
		comm.CreateAt = time.Time{}
		comm.UpdateAt = time.Time{}
		comm.TotalLike = 0
		comm.Num = 0
		comm.ID = 0
		commentPool.Put(comm)
	}()

	tx.Model(comm).Select("id", "content", "total_like", "time", "num", "like", "user_id", "rpid").
		Where("content = ?", c.Content).Attrs(Comment{
		Content:   c.Content,
		Time:      c.Time,
		TotalLike: 0,
		Num:       0,
		Like:      c.Like,
		UserID:    c.UserID,
		Rpid:      c.Rpid,
	}).FirstOrCreate(comm)

	cache.GetCache().Check.Increment("check", fmt.Sprint(comm.ID), check.HashSet(comm.Content))
	cache.GetCache().Content.Increment("content", fmt.Sprint(comm.ID), comm.Content)

	c.CommentID = comm.ID

	if comm.Rpid == c.Rpid {
		return nil
	}

	if comm.Time > c.Time {
		comm.Rpid = c.Rpid
		comm.Time = c.Time
		comm.Like = c.Like
		comm.UserID = c.UserID
	}
	comm.Num++
	comm.TotalLike += c.Like

	tx.Model(comm).Select("Rpid", "TotalLike", "Time", "Num", "Like", "UserID").Updates(comm)
	return nil
}

// Commentators 批量插入使用
type Commentators []*Commentator

// GetClauses 插入时冲突解决方法
func (Commentators) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		DoNothing: true,
	}
}

// GetModels 查询时返回的切片
func (Commentators) GetModels() interface{} {
	return nil
}

// TableName 表名称
func (Commentators) TableName() string {
	return "commentator"
}
