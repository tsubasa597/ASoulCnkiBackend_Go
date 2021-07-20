package db

type Dynamic struct {
	Model
	DynamicID int64 `json:"id" gorm:"primaryKey;column:id;uniqueIndex"`
	Type      int8  `json:"type" gorm:"column:type"`
	StartTime int64 `json:"start_time" gorm:"column:start_time"`
	EndTime   int64 `json:"end_time" gorm:"column:end_time"` // 最后评论的时间
}

var _ Modeler = (*Dynamic)(nil)

func (Dynamic) getModels() interface{} {
	return &[]Dynamic{}
}

func (Dynamic) TableName() string {
	return "dynamic"
}
