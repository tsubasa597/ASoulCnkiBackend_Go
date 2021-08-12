package entry

type User struct {
	Model
	UID             int64  `json:"uid" gorm:"uniqueIndex,column:uid"`
	Name            string `json:"name" gorm:"column:name"`
	LastDynamicTime int32  `json:"dynamic_time" gorm:"column:dynamic_time"`
}

var _ Modeler = (*User)(nil)

func (User) GetModels() interface{} {
	return &[]User{}
}

func (User) TableName() string {
	return "user"
}
