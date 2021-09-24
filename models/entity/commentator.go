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

type Commentator struct {
	Model
	UID       int64  `json:"uid" gorm:"column:uid"`
	UName     string `json:"uname" gorm:"column:uname"`
	Like      uint32 `json:"like" gorm:"column:like"`
	Time      int64  `json:"-" gorm:"column:time"`
	Content   string `json:"-" gorm:"-"`
	DynamicID uint64 `json:"-" gorm:"index:idx_dynamic_id"`
	UserID    uint64 `json:"-" grom:"-"`
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

func (Commentator) TableName() string {
	return "commentator"
}

func (Commentator) GetModels() interface{} {
	return &[]Commentator{}
}

func (c *Commentator) BeforeCreate(tx *gorm.DB) error {
	comm := commentPool.Get().(*Comment)
	defer func() {
		comm.ID = 0
		comm.CreateAt = time.Time{}
		comm.UpdateAt = time.Time{}
		comm.TotalLike = 0
		comm.Num = 0
		comm.Like = 0
		comm.UserID = 0
		commentPool.Put(comm)
	}()

	tx.Model(comm).Select("id", "content", "total_like", "time", "num", "like", "user_id").
		Where("content = ?", c.Content).Attrs(Comment{
		Model: Model{
			ID: c.ID,
		},
		Content:   c.Content,
		Time:      c.Time,
		TotalLike: 0,
		Num:       0,
		Like:      c.Like,
		UserID:    c.UserID,
	}).FirstOrCreate(comm)

	cache.GetCache().Check.Increment(fmt.Sprint(comm.ID), check.HashSet(comm.Content))
	cache.GetCache().Content.Set(fmt.Sprint(comm.ID), comm.Content)

	if comm.ID == c.ID {
		return nil
	}

	if comm.Time > c.Time {
		comm.ID = c.ID
		comm.Time = c.Time
		comm.Like = c.Like
		comm.UserID = c.UserID
	}
	comm.Num++
	comm.TotalLike += c.Like

	tx.Model(comm).Select("ID", "TotalLike", "Time", "Num", "Like", "UserID").Updates(comm)
	return nil
}

type Commentators []*Commentator

func (Commentators) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		DoNothing: true,
	}
}

func (Commentators) GetModels() interface{} {
	return nil
}

func (Commentators) TableName() string {
	return "commentator"
}
