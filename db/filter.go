package db

import (
	"fmt"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
	"gorm.io/gorm"
)

type Param struct {
	Page  int
	Order string
	Field []string
	Query string
	Args  []interface{}
}

type Time uint8

type Sort uint8

const (
	All      Time = 0
	OneWeek  Time = 1
	ThreeDay Time = 2

	TotalLikeSort Sort = 1
	LikeSort      Sort = 2
	NumSort       Sort = 3
)

func NewFilter(t Time, users []entry.User, sort Sort) (p Param) {
	p.Query = "time > ?"
	switch t {
	case OneWeek:
		p.Args = []interface{}{time.Since(time.Now().AddDate(0, 0, 7))}
	case ThreeDay:
		p.Args = []interface{}{time.Since(time.Now().AddDate(0, 0, 7))}
	default:
		p.Args = []interface{}{0}
	}

	switch sort {
	case TotalLikeSort:
		p.Order = "total_like desc"
	case LikeSort:
		p.Order = "like desc"
	case NumSort:
		p.Order = "num desc"
	}

	if len(users) < 1 {
		return
	}

	p.Query += " and user_id in ?"
	uids := make([]int64, len(users))
	for _, user := range users {
		uids = append(uids, int64(user.ID))
	}
	p.Args = append(p.Args, uids)

	return
}

func FilterTime(timestamp int32) Param {
	return Param{
		Query: "time > ?",
		Args:  []interface{}{timestamp},
		Order: "like desc",
	}
}

func FilterUser(users []entry.User) Param {
	p := Param{}
	if len(users) == 0 {
		return p
	}

	for _, user := range users {
		p.Args = append(p.Args, fmt.Sprint(user.UID))
	}
	p.Query = "uid in ?"
	return p
}

func filter(param Param) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		size := conf.Size
		if param.Page < 0 {
			param.Page = 2
			size = -1
		}

		// if param.Preload == nil {
		// 	return db.Select(param.Field).Where(param.Query, param.Args...).
		// 		Offset((param.Page - 1) * size).Limit(size).Order(param.Order)
		// }

		db = db.Select(param.Field)

		// for i := range param.PreloadField {
		// 	db = db.Preload(param.PreloadField[i])
		// }

		return db.Where(param.Query, param.Args...).
			Offset((param.Page - 1) * size).Limit(size).Order(param.Order)
	}
}
