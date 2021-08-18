package entry

import (
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

type Commentator struct {
	Model
	UID       int64  `json:"uid" gorm:"column:uid"`
	UName     string `json:"uname" gorm:"column:uname"`
	Rpid      int64  `json:"rpid" gorm:"column:rpid"`
	Like      uint32 `json:"like" gorm:"column:like"`
	Time      int64  `json:"-" gorm:"column:time"`
	Content   string `json:"-" gorm:"-"`
	UserID    uint64 `json:"-" gorm:"-"`
	CommentID uint64 `json:"-"`
	DynamicID uint64 `json:"-"`
}

var (
	_           Modeler                         = (*Commentator)(nil)
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
		commentPool.Put(comm)
	}()

	tx.Model(comm).Select("id", "like", "content", "total_like", "time", "num", "rpid", "user_id").
		Where("content = ?", c.Content).Attrs(Comment{
		Content:   c.Content,
		Like:      c.Like,
		Time:      c.Time,
		Rpid:      c.Rpid,
		TotalLike: 0,
		Num:       0,
		UserID:    c.UserID,
	}).FirstOrCreate(comm)

	c.CommentID = comm.ID

	if comm.Rpid == c.Rpid {
		return nil
	}

	if comm.Time > c.Time {
		comm.UserID = c.UserID
		comm.Rpid = c.Rpid
		comm.Time = c.Time
		comm.Like = c.Like
	}
	comm.Num += 1
	comm.TotalLike += comm.Like

	tx.Model(comm).Select("Rpid", "TotalLike", "Time", "Num", "Like", "UserID").
		Updates(comm)
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
