package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"

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
			&gorm.Config{})
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
	db.db.AutoMigrate(&entry.Comment{}, &entry.Dynamic{}, &entry.Article{})
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

const (
	ErrNotFound      = "数据不存在"
	ErrHasNoDataBase = "没有指定数据库"
)
