package model

import (
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
	"gorm.io/gorm"
)

// Param 查询筛选条件
type Param struct {
	Page  int
	Size  int
	Order string
	Field []string
	Query string
	Args  []interface{}
}

func filter(param Param) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if param.Size < 1 {
			param.Size = setting.Size
		}

		if param.Page < 0 {
			param.Page = 2
			param.Size = -1
		}

		db = db.Select(param.Field)

		return db.Where(param.Query, param.Args...).
			Offset((param.Page - 1) * param.Size).Limit(param.Size).Order(param.Order)
	}
}
