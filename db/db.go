package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
)

var (
	db    *gorm.DB
	mutex sync.Mutex
)

type Modeler interface {
	getModels() interface{}
}

type Model struct {
	ID       uint64    `json:"-" gorm:"primaryKey;autoIncrement"`
	CreateAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"-" gorm:"autoUpdateTime"`
}

func init() {
	var err error
	if conf.SQL == "sqllite" {
		db, err = gorm.Open(sqlite.Open("comment.db"), &gorm.Config{})
	} else if conf.SQL == "mysql" {
		db, err = gorm.Open(
			mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				conf.User, conf.Password, conf.Host, conf.DBName)),
			&gorm.Config{})
	} else {
		panic(ErrHasNoDataBase)
	}

	if err != nil {
		panic(err)
	}

	db.Use(dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(conf.MaxOpen).
		SetMaxOpenConns(conf.MaxConn),
	)

	if conf.RunMode == "debug" {
		db.Debug()
	}

	migrateTable()
}

func migrateTable() {
	if !db.Migrator().HasTable(&User{}) {
		db.AutoMigrate(&User{})
		fmt.Println(Add(&User{
			UID:             351609538,
			LastDynamicTime: 1606403616, // 1627381252
		}))
		fmt.Println(Add(&User{
			UID:             672328094,
			LastDynamicTime: 1606133780, // 1627381148
		}))
		fmt.Println(Add(&User{
			UID:             672353429,
			LastDynamicTime: 1606403340,
		}))
		fmt.Println(Add(&User{
			UID:             672346917,
			LastDynamicTime: 1606403478,
		}))
		fmt.Println(Add(&User{
			UID:             672342685,
			LastDynamicTime: 1606403225,
		}))
	}
	db.AutoMigrate(&Comment{}, &Dynamic{})
}

func Get(model Modeler) interface{} {
	mutex.Lock()
	defer mutex.Unlock()

	models := model.getModels()
	db.Find(models)

	return models
}

func Find(model Modeler) error {
	mutex.Lock()
	defer mutex.Unlock()

	if db.Where(model).Find(model).RowsAffected == 0 {
		return fmt.Errorf(ErrNotFound)
	}
	return db.Error
}

func Add(model Modeler) error {
	mutex.Lock()
	defer mutex.Unlock()

	if conf.SQL == "mysql" {
		return db.Clauses(clause.Insert{
			Modifier: "IGNORE",
		}).Create(model).Error
	} else if conf.SQL == "sqllite" {
		return db.Clauses(clause.Insert{
			Modifier: "OR IGNORE",
		}).Create(model).Error
	}
	return fmt.Errorf(ErrHasNoDataBase)
}

func Update(model Modeler) error {
	mutex.Lock()
	defer mutex.Unlock()

	return db.Updates(model).Error
}

func Delete(model Modeler) error {
	mutex.Lock()
	defer mutex.Unlock()

	return db.Delete(model).Error
}

const (
	ErrNotFound      = "不存在"
	ErrHasNoDataBase = "没有指定数据库"
)
