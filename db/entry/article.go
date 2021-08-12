package entry

type Article struct {
	Model
	Like        uint32 `json:"like" gorm:"column:like"`
	TotalLike   uint32 `json:"total_like" gorm:"column:total_like"`
	Time        int64  `json:"comment_time" gorm:"column:time"`
	Num         uint32 `json:"num" gorm:"column:num"`
	CommentText string `json:"commnet" gorm:"column:comment_text;uniqueIndex"`
	UserID      uint64
	CommentID   uint64
	Comment     *Comment `gorm:"foreignKey:CommentID"`
}

var _ Modeler = (*Article)(nil)

func (Article) GetModels() interface{} {
	return &[]Article{}
}

func (Article) TableName() string {
	return "article"
}
