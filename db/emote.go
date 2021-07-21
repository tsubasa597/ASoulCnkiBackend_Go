package db

type Emote struct {
	Model
	EmoteID   int32  `json:"id" gorm:"column:emote_id;uniqueIndex"`
	EmoteText string `json:"emote" gorm:"column:emote;uniqueIndex"`
}

var _ Modeler = (*Emote)(nil)

func (Emote) getModels() interface{} {
	return &[]Emote{}
}

func (Emote) TableName() string {
	return "emote"
}
