package db

import (
	"fmt"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
)

var (
	db *gorm.DB
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

	if !db.Migrator().HasTable(&User{}) {
		db.AutoMigrate(&User{})
		fmt.Println(Add(&User{
			UID:             672328094,
			LastDynamicTime: 1626500626, // 99164512 3 M

		}))
	}
	db.AutoMigrate(&Comment{}, &Dynamic{}, &Emote{})

}

func MigrateAll(models []Modeler) {
	for _, v := range models {
		db.AutoMigrate(v)
	}
}

func Get(model Modeler) interface{} {
	models := model.getModels()
	db.Find(models)
	return models
}

func Find(model Modeler) error {
	if db.Where(model).Find(model).RowsAffected == 0 {
		return fmt.Errorf(ErrNotFound)
	}
	return db.Error
}

func Add(model Modeler) error {
	if conf.SQL == "mysql" {
		return db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(model).Error
	} else if conf.SQL == "sqllite" {
		return db.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(model).Error
	}
	return fmt.Errorf(ErrHasNoDataBase)
}

func Update(model Modeler) error {
	return db.Updates(model).Error
}

func Delete(model Modeler) error {
	return db.Delete(model).Error
}

const (
	ErrNotFound      = "不存在"
	ErrHasNoDataBase = "没有指定数据库"
)
