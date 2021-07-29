package db

type User struct {
	Model
	UID             int64  `json:"uid" gorm:"column:uid"`
	Name            string `json:"name" gorm:"column:name"`
	LastDynamicTime int32  `json:"dynamic_time" gorm:"column:dynamic_time"`
}

var _ Modeler = (*User)(nil)

func (User) getModels() interface{} {
	return &[]User{}
}

func (User) TableName() string {
	return "user"
}
