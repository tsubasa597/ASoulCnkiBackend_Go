package entity

type User struct {
	Model
	UID             int64  `json:"uid" gorm:"column:uid;uniqueIndex"`
	Name            string `json:"name" gorm:"column:name"`
	LastDynamicTime int32  `json:"dynamic_time" gorm:"column:dynamic_time"`
}

var _ Entity = (*User)(nil)

func (User) GetModels() interface{} {
	return &[]User{}
}

func (User) TableName() string {
	return "user"
}
