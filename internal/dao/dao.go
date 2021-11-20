package dao

import (
	"fmt"
	"time"

	e "github.com/tsubasa597/ASoulCnkiBackend/internal/ecode"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	_db *gorm.DB
)

// Setup 初始化数据库
func Setup() error {
	var err error
	log, err := config.NewLogFile("database")
	if err != nil {
		return err
	}

	switch config.SQL {
	case "mysql":
		_db, err = gorm.Open(
			mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				config.User, config.Password, config.Host, config.DBName)),
			&gorm.Config{
				Logger: logger.New(log, logger.Config{
					SlowThreshold: time.Second * 5,
					LogLevel:      logger.Error,
				}),
				PrepareStmt: true,
			})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf(e.ErrHasNoDataBase)
	}

	if err := _db.Use(dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(config.MaxOpen).
		SetMaxOpenConns(config.MaxConn),
	); err != nil {
		return err
	}

	if config.RunMode == "debug" {
		_db = _db.Debug()
	}

	if err := migrateTable(); err != nil {
		return err
	}

	return nil
}

func migrateTable() error {
	if !_db.Migrator().HasTable(&entity.User{}) {
		if err := _db.AutoMigrate(&entity.User{}); err != nil {
			return err
		}

		err := _db.Create([]entity.User{
			{
				UID:             351609538,
				LastDynamicTime: 1606403616,
			},
			{
				UID:             672328094,
				LastDynamicTime: 1606133780,
			},
			{
				UID:             672353429,
				LastDynamicTime: 1606403340,
			},
			{
				UID:             672346917,
				LastDynamicTime: 1606403478,
			},
			{
				UID:             672342685,
				LastDynamicTime: 1606403225,
			},
			{
				UID:             703007996,
				LastDynamicTime: 1606118482,
			},
		}).Error
		if err != nil {
			return err
		}
	}

	return _db.AutoMigrate(&entity.Comment{}, &entity.Dynamic{}, &entity.Commentator{}, &entity.User{})
}

// Get 获取所有数据
func Get(model entity.Entity) (interface{}, error) {
	models := model.GetModels()
	_db.Find(models)

	return models, _db.Error
}

// Add 添加数据
func Add(model entity.Entity, batchSize int) error {
	return _db.Clauses(model.GetClauses()).CreateInBatches(model, batchSize).Error
}
