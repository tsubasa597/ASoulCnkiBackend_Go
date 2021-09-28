package entity

import (
	"time"

	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Entity 数据库表接口
type Entity interface {
	schema.Tabler
	GetModels() interface{}
	GetClauses() clause.OnConflict
}

// Model 公共结构
type Model struct {
	ID       uint64    `json:"-" gorm:"primaryKey;autoIncrement;autoIncrementIncrement:1"`
	CreateAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"-" gorm:"autoUpdateTime"`
}

// GetClauses 插入时冲突解决方法
func (Model) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		DoNothing: true,
	}
}
