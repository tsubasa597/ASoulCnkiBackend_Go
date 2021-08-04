package db

type Dynamic struct {
	Model
	DynamicID int64 `json:"dynamic_id" gorm:"primaryKey;column:dynamic_id;uniqueIndex"`
	UID       int64 `json:"uid" gorm:"column:uid"`
	RID       int64 `json:"rid" gorm:"column:rid;uniqueIndex"`
	Type      uint8 `json:"type" gorm:"column:type"`
	Time      int32 `json:"time" gorm:"column:time"`
	Updated   bool  `json:"is_update" gorm:"column:is_update"`
}

var _ Modeler = (*Dynamic)(nil)

func (Dynamic) getModels() interface{} {
	return &[]Dynamic{}
}

func (Dynamic) TableName() string {
	return "dynamic"
}

func (dynamic Dynamic) Find(params []interface{}) interface{} {

	models := dynamic.getModels()
	db.Order("time asc").Find(models, params...)
	return models
}
