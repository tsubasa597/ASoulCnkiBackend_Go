package entry

import (
	"sync"

	"gorm.io/gorm"
)

type Dynamic struct {
	Model
	RID     int64 `json:"rid" gorm:"column:rid;uniqueIndex"`
	Type    uint8 `json:"type" gorm:"column:type"`
	Time    int32 `json:"time" gorm:"column:time"`
	Updated bool  `json:"is_update" gorm:"column:is_update"`
	UserID  uint64
}

var (
	_        Modeler   = (*Dynamic)(nil)
	userPool sync.Pool = sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}
)

func (d *Dynamic) AfterCreate(tx *gorm.DB) error {
	user := userPool.Get().(*User)
	defer func() {
		user.ID = 0
		user.LastDynamicTime = 0
		userPool.Put(user)
	}()

	if err := tx.Model(user).Select("id", "dynamic_time").Where("id", d.UserID).Find(user).Error; err != nil {
		return err
	}

	user.LastDynamicTime = d.Time
	return tx.Model(&user).Select("dynamic_time").Updates(user).Error
}

func (Dynamic) GetModels() interface{} {
	return &[]Dynamic{}
}

func (Dynamic) TableName() string {
	return "dynamic"
}
