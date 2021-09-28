package entity

// User 用户表
type User struct {
	Model
	UID             int64  `json:"uid" gorm:"column:uid;uniqueIndex"`
	Name            string `json:"name" gorm:"column:name"`
	LastDynamicTime int32  `json:"dynamic_time" gorm:"column:dynamic_time"`
}

var _ Entity = (*User)(nil)

// GetModels 查询时返回的切片
func (User) GetModels() interface{} {
	return &[]User{}
}

// TableName 表名称
func (User) TableName() string {
	return "user"
}
