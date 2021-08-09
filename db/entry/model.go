package entry

import "time"

type Modeler interface {
	GetModels() interface{}
}

type Model struct {
	ID       uint64    `json:"-" gorm:"primaryKey;autoIncrement"`
	CreateAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"-" gorm:"autoUpdateTime"`
}
