package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
	"github.com/tsubasa597/ASoulCnkiBackend/db/vo"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type DB struct {
	db    *gorm.DB
	mutex *sync.Mutex
}

func New() (*DB, error) {
	var (
		db = &DB{
			mutex: &sync.Mutex{},
		}
		err error
	)

	if conf.SQL == "sqlite" {
		db.db, err = gorm.Open(sqlite.Open("comment.db"), &gorm.Config{})
	} else if conf.SQL == "mysql" {
		db.db, err = gorm.Open(
			mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				conf.User, conf.Password, conf.Host, conf.DBName)),
			&gorm.Config{
				PrepareStmt: true,
			})
	} else {
		return nil, fmt.Errorf(ErrHasNoDataBase)
	}

	if err != nil {
		return nil, err
	}

	db.db.Use(dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(conf.MaxOpen).
		SetMaxOpenConns(conf.MaxConn),
	)

	if conf.RunMode == "debug" {
		db.db = db.db.Debug()
	}

	db.migrateTable()

	return db, nil
}

func (db DB) migrateTable() {
	if !db.db.Migrator().HasTable(&entry.User{}) {
		db.db.AutoMigrate(&entry.User{})
		if err := db.Add(&entry.User{
			UID:             351609538,
			LastDynamicTime: 1606403616, // 1627381252
		}); err != nil {
			panic(err)
		}
		if err := db.Add(&entry.User{
			UID:             672328094,
			LastDynamicTime: 1606133780, // 1627381148
		}); err != nil {
			panic(err)
		}
		if err := db.Add(&entry.User{
			UID:             672353429,
			LastDynamicTime: 1606403340,
		}); err != nil {
			panic(err)
		}
		if err := db.Add(&entry.User{
			UID:             672346917,
			LastDynamicTime: 1606403478,
		}); err != nil {
			panic(err)
		}
		if err := db.Add(&entry.User{
			UID:             672342685,
			LastDynamicTime: 1606403225,
		}); err != nil {
			panic(err)
		}
	}
	db.db.AutoMigrate(&entry.Comment{}, &entry.Dynamic{}, &entry.Commentator{})
}

func (db DB) Get(model entry.Modeler) (interface{}, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	models := model.GetModels()
	db.db.Find(models)

	return models, db.db.Error
}

func (db DB) Find(model entry.Modeler, param Param) (interface{}, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	models := model.GetModels()

	if db.db.Scopes(filter(param)).Find(models).RowsAffected == 0 {
		return models, fmt.Errorf(ErrNotFound)
	}
	return models, db.db.Error
}

func (db DB) Add(model entry.Modeler) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.db.Clauses(model.GetClauses()).Create(model).Error
}

func (db DB) Update(model entry.Modeler, param Param) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.db.Model(model).Select(param.Field).Updates(model).Error
}

func (db DB) Delete(model entry.Modeler) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	return db.db.Delete(model).Error
}

func (db DB) Rank(page, size int, time, order string, uids ...string) (replys []vo.Reply, err error) {
	tx := db.getReply(page, size)

	for _, uid := range uids {
		tx.Where("user.uid = ?", uid)
	}

	if time >= "0" {
		tx.Where("comment.time > ?", time)
	}

	if order != "" {
		tx.Order("comment." + order)
	}

	tx.Find(&replys)
	return replys, tx.Error
}

func (db DB) Check(rpid string) (vo.Reply, error) {
	reply := vo.Reply{}

	tx := db.getReply(2, -1)
	tx.Where("comment.rpid = ?", rpid).First(&reply)

	return reply, tx.Error
}

type CommentCache struct {
	Rpid    int64
	Content string
}

func (db DB) GetContent(rpid string) (commentCache []CommentCache) {
	db.db.Model(&entry.Commentator{}).Select("commentator.rpid, comment.content").
		Joins("left join comment on commentator.comment_id = comment.id").
		Where("commentator.rpid > ?", rpid).
		Find(&commentCache)

	return
}

func (db DB) getReply(page, size int) gorm.DB {
	return *db.db.Model(&entry.Commentator{}).
		Select("dynamic.type, dynamic.rid as rid, user.uid as uuid, commentator.rpid, commentator.uid as uid, commentator.time, commentator.uname as name, comment.content, commentator.like, comment.rpid as origin_rpid, comment.num, comment.total_like").
		Joins("left join comment on comment.id = commentator.comment_id").
		Joins("left join dynamic on dynamic.id = commentator.dynamic_id").
		Joins("left join user on user.id = dynamic.user_id").
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
			param.Size = conf.Size
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

const (
	ErrNotFound      = "数据不存在"
	ErrHasNoDataBase = "没有指定数据库"
)
