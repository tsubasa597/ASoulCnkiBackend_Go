package entity

import (
	"sort"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

// Dynamic 动态表
type Dynamic struct {
	Model
	RID     int64  `json:"rid" gorm:"column:rid;uniqueIndex"`
	Type    uint8  `json:"type" gorm:"column:type"`
	Time    int32  `json:"time" gorm:"column:time"`
	Updated bool   `json:"is_update" gorm:"column:is_update"`
	Name    string `json:"-" gorm:"-"`
	UserID  uint64 `json:"-" gorm:"index:idx_user_id"`
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

// AfterCreate 插入之后更新用户表中数据
func (d *Dynamic) AfterCreate(tx *gorm.DB) error {
	user := userPool.Get().(*User)
	defer userPool.Put(user)

	user.LastDynamicTime = d.Time
	user.Name = d.Name

	return tx.Model(&user).Select("name", "dynamic_time").Where("id = ?", d.UserID).Updates(user).Error
}

// GetModels 查询时返回的切片
func (Dynamic) GetModels() interface{} {
	return &[]Dynamic{}
}

// TableName 表名称
func (Dynamic) TableName() string {
	return "dynamic"
}

// Dynamics 批量插入使用
type Dynamics []*Dynamic

func (d Dynamics) Len() int           { return len(d) }
func (d Dynamics) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Dynamics) Less(i, j int) bool { return d[i].Time < d[j].Time }

// GetModels 查询时返回的切片
func (Dynamics) GetModels() interface{} {
	return &[]Dynamic{}
}

// TableName 表名称
func (Dynamics) TableName() string {
	return "dynamic"
}

// GetClauses 插入时冲突解决方法
func (Dynamics) GetClauses() clause.OnConflict {
	return clause.OnConflict{
		DoNothing: true,
	}
}
