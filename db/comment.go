package db

type Comment struct {
	Model
	UID       int64  `json:"uid" gorm:"column:uid"`
	UName     string `json:"uname" gorm:"column:uname"`
	Comment   string `json:"comment" gorm:"column:comment;uniqueIndex"`
	CommentID int64  `json:"comment_id" gorm:"column:comment_id"`
	Time      int64  `json:"comment_time" gorm:"column:time"`
	// Like      uint32 `json:"like" gorm:"column:like"`
}

var _ Modeler = (*Comment)(nil)

func (Comment) getModels() interface{} {
	return &[]Comment{}
}

func (Comment) TableName() string {
	return "comment"
}

type Comments []*Comment

var _ Modeler = (*Comments)(nil)

func (Comments) getModels() interface{} {
	return &[]Comment{}
}

func (Comments) TableName() string {
	return "comment"
}
