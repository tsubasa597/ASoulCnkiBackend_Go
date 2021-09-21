package model

import (
	"fmt"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/models/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/models/vo"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/e"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
)

var (
	db *gorm.DB
)

func Setup() {
	var err error

	switch setting.SQL {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open("comment.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	case "mysql":
		db, err = gorm.Open(
			mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				setting.User, setting.Password, setting.Host, setting.DBName)),
			&gorm.Config{
				PrepareStmt: true,
			})
		if err != nil {
			panic(err)
		}
	default:
		panic(e.ErrHasNoDataBase)
	}

	db.Use(dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(setting.MaxOpen).
		SetMaxOpenConns(setting.MaxConn),
	)

	if setting.RunMode == "debug" {
		db = db.Debug()
	}

	migrateTable()
}

func migrateTable() {
	if !db.Migrator().HasTable(&entity.User{}) {
		db.AutoMigrate(&entity.User{})
		if err := Add(&entity.User{
			UID:             351609538,
			LastDynamicTime: 1606403616, // 1627381252
		}); err != nil {
			panic(err)
		}
		if err := Add(&entity.User{
			UID:             672328094,
			LastDynamicTime: 1606133780, // 1627381148
		}); err != nil {
			panic(err)
		}
		if err := Add(&entity.User{
			UID:             672353429,
			LastDynamicTime: 1606403340,
		}); err != nil {
			panic(err)
		}
		if err := Add(&entity.User{
			UID:             672346917,
			LastDynamicTime: 1606403478,
		}); err != nil {
			panic(err)
		}
		if err := Add(&entity.User{
			UID:             672342685,
			LastDynamicTime: 1606403225,
		}); err != nil {
			panic(err)
		}
	}
	db.AutoMigrate(&entity.Comment{}, &entity.Dynamic{}, &entity.Commentator{})
}

func Get(model entity.Entity) (interface{}, error) {
	models := model.GetModels()
	db.Find(models)

	return models, db.Error
}

func Find(model entity.Entity, param Param) (interface{}, error) {
	models := model.GetModels()

	if db.Scopes(filter(param)).Find(models).RowsAffected == 0 {
		return models, fmt.Errorf(e.ErrNotFound)
	}
	return models, db.Error
}

func Add(model entity.Entity) error {
	return db.Clauses(model.GetClauses()).Create(model).Error
}

func Update(model entity.Entity, param Param) error {
	return db.Model(model).Select(param.Field).Updates(model).Error
}

func Delete(model entity.Entity) error {
	return db.Delete(model).Error
}

func Rank(page, size int, time, order string, uids ...string) (replys []vo.Reply, err error) {
	tx := getReply(page, size)

	for _, uid := range uids {
		tx.Where("user.uid = ?", uid)
	}

	if time >= "0" {
		tx.Where("comment.time > ?", time)
	}

	if order != "" {
		tx.Order(clause.OrderByColumn{
			Column: clause.Column{Name: "comment." + order},
			Desc:   true,
		})
	}

	tx.Find(&replys)
	return replys, tx.Error
}

func Check(rpid string) (vo.Reply, error) {
	reply := vo.Reply{}

	tx := getReply(2, -1)
	tx.Where("comment.rpid = ?", rpid).First(&reply)

	return reply, tx.Error
}

type CommentCache struct {
	Rpid    int64
	Content string
}

func GetContent(rpid string) (commentCache []CommentCache) {
	db.Model(&entity.Commentator{}).Select("commentator.rpid, comment.content").
		Joins("left join comment comment on commentator.comment_id = comment.id").
		Where("commentator.rpid > ?", rpid).
		Order("commentator.rpid asc").
		Find(&commentCache)

	return
}

func getReply(page, size int) gorm.DB {
	return *db.Model(&entity.Comment{}).
		Select("dynamic.type, dynamic.rid as rid, user.uid as uuid, commentator.rpid, commentator.uid as uid, commentator.time, commentator.uname as name, comment.content, commentator.like, comment.rpid as origin_rpid, comment.num, comment.total_like").
		Joins("left join commentator commentator on comment.rpid = commentator.rpid").
		Joins("left join dynamic dynamic on dynamic.rid = commentator.dynamic_id").
		Joins("left join user user on user.id = dynamic.user_id").
		Limit(size).Offset(size * (page - 1))
}

type Param struct {
	Page  int
	Size  int
	Order string
	Field []string
	Query string
	Args  []interface{}
}

func filter(param Param) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if param.Size < 1 {
			param.Size = setting.Size
		}

		if param.Page < 0 {
			param.Page = 2
			param.Size = -1
		}

		db = db.Select(param.Field)

		return db.Where(param.Query, param.Args...).
			Offset((param.Page - 1) * param.Size).Limit(param.Size).Order(param.Order)
	}
}
