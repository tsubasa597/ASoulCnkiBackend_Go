package entry

import (
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

type Dynamic struct {
	Model
	RID     int64  `json:"rid" gorm:"column:rid;uniqueIndex"`
	Type    uint8  `json:"type" gorm:"column:type"`
	Time    int32  `json:"time" gorm:"column:time"`
	Updated bool   `json:"is_update" gorm:"column:is_update"`
	Name    string `json:"-" gorm:"-"`
	UserID  uint64
}

var (
	_        Modeler                        = (*Dynamic)(nil)
	_        callbacks.AfterCreateInterface = (*Dynamic)(nil)
	userPool sync.Pool                      = sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}
)

func (d *Dynamic) AfterCreate(tx *gorm.DB) error {
	user := userPool.Get().(*User)
	defer func() {
		user.Name = ""
		user.LastDynamicTime = 0
		userPool.Put(user)
	}()

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
