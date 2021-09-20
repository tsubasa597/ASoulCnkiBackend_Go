package entity

import (
	"time"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Entity interface {
	schema.Tabler
	GetModels() interface{}
	GetClauses() clause.OnConflict
}

type Model struct {
	ID       uint64    `json:"-" gorm:"primaryKey;autoIncrement;autoIncrementIncrement:1"`
	CreateAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"-" gorm:"autoUpdateTime"`
}

func (Model) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		DoNothing: true,
	}
}
