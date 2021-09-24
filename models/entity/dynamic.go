package entity

import (
	"sort"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

type Dynamic struct {
	Model
	RID     int64  `json:"rid" gorm:"column:rid;uniqueIndex"`
	Type    uint8  `json:"type" gorm:"column:type"`
	Time    int32  `json:"time" gorm:"column:time"`
	Updated bool   `json:"is_update" gorm:"column:is_update"`
	Name    string `json:"-" gorm:"-"`
	UserID  uint64 `json:"-" gorm:"uniqueIndex"`
}

var (
	_, _     Entity                         = (*Dynamic)(nil), (*Dynamics)(nil)
	_        callbacks.AfterCreateInterface = (*Dynamic)(nil)
	_        sort.Interface                 = (*Dynamics)(nil)
	userPool sync.Pool                      = sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}
)

func (d *Dynamic) AfterCreate(tx *gorm.DB) error {
	user := userPool.Get().(*User)
	defer userPool.Put(user)

	user.LastDynamicTime = d.Time
	user.Name = d.Name

	return tx.Model(&user).Select("name", "dynamic_time").Where("id = ?", d.UserID).Updates(user).Error
}

func (Dynamic) GetModels() interface{} {
	return &[]Dynamic{}
}

func (Dynamic) TableName() string {
	return "dynamic"
}

type Dynamics []*Dynamic

func (d Dynamics) Len() int           { return len(d) }
func (d Dynamics) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Dynamics) Less(i, j int) bool { return d[i].Time < d[j].Time }

func (Dynamics) GetModels() interface{} {
	return &[]Dynamic{}
}

func (Dynamics) TableName() string {
	return "dynamic"
}

func (Dynamics) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		DoNothing: true,
	}
}
