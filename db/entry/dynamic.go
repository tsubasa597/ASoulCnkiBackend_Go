package entry

type Dynamic struct {
	Model
	RID     int64 `json:"rid" gorm:"column:rid;uniqueIndex"`
	Type    uint8 `json:"type" gorm:"column:type"`
	Time    int32 `json:"time" gorm:"column:time"`
	Updated bool  `json:"is_update" gorm:"column:is_update"`
	UserID  uint64
}

var _ Modeler = (*Dynamic)(nil)

func (Dynamic) GetModels() interface{} {
	return &[]Dynamic{}
}

func (Dynamic) TableName() string {
	return "dynamic"
}
