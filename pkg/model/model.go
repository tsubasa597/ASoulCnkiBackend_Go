package model

import (
	"fmt"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/models/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/e"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	db *gorm.DB
)

// Setup 初始化数据库
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
			LastDynamicTime: 1606403616,
		}); err != nil {
			panic(err)
		}
		if err := Add(&entity.User{
			UID:             672328094,
			LastDynamicTime: 1606133780,
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
	db.AutoMigrate(&entity.Comment{}, &entity.Dynamic{}, &entity.Commentator{}, &entity.User{})
}

// Get 获取所有数据
func Get(model entity.Entity) (interface{}, error) {
	models := model.GetModels()
	db.Find(models)

	return models, db.Error
}

// Find 条件查询
func Find(model entity.Entity, param Param) (interface{}, error) {
	models := model.GetModels()

	if db.Scopes(filter(param)).Find(models).RowsAffected == 0 {
		return models, fmt.Errorf(e.ErrNotFound)
	}
	return models, db.Error
}

// Add 添加数据
func Add(model entity.Entity) error {
	return db.Clauses(model.GetClauses()).Create(model).Error
}

// Update 更新数据
func Update(model entity.Entity, param Param) error {
	return db.Model(model).Select(param.Field).Updates(model).Error
}

// Delete 删除数据
func Delete(model entity.Entity) error {
	return db.Delete(model).Error
}
