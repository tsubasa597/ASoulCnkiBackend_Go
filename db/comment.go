package db

type Comment struct {
	Model
	UID       int64  `json:"uid" gorm:"column:uid"`
	UName     string `json:"uname" gorm:"column:uname"`
	Comment   string `json:"comment" gorm:"column:comment;uniqueIndex"`
	DynamicID int64  `json:"dynamic_id" gorm:"column:dynamic_id"`
	Time      int64  `json:"comment_time" gorm:"column:time"`
}

var _ Modeler = (*Comment)(nil)

func (Comment) getModels() interface{} {
	return &[]Comment{}
}

func (Comment) TableName() string {
	return "comment"
}
