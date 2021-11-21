package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/dao"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/server"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/service/listen"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

var (
	_log *logrus.Logger
)

func main() {
	if err := setup(); err != nil {
		_log.Fatal(err)
	}

	// 退出后处理
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := server.HttpServer.Shutdown(ctx); err != nil {
		_log.Error(err)
	}

	if err := cache.GetInstance().Stop(); err != "" {
		_log.Error(err)
	}

	if err := listen.GetInstance().Stop(); err != nil {
		_log.Error(err)
	}
}

func setup() error {
	if err := loadConf(); err != nil {
		return err
	}

	var err error
	// 系统日志写入
	_log, err = config.NewLogFile("system")
	if err != nil {
		return err
	}

	// 初始化数据库
	if err := dao.Setup(); err != nil {
		return err
	}

	// 初始化缓存
	if err := cache.Setup(); err != nil {
		return err
	}

	// 初始化自动更新
	if err = listen.Setup(); err != nil {
		return err
	}

	// 初始化 http 监听
	if err := server.NewHttpServer(); err != nil {
		return err
	}

	go func() {
		if err := server.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			_log.Error(err)
		}
	}()

	return nil
}

func loadConf() error {
	cfg, err := ini.Load("./conf/conf.ini")
	if err != nil {
		return fmt.Errorf("Fail to parse 'conf/app.ini': %v", err)
	}
	config.RunMode = cfg.Section("web").Key("RUN_MODE").MustString("debug")
	config.Port = cfg.Section("web").Key("PORT").MustInt(8000)
	config.Size = cfg.Section("web").Key("SIZE").MustInt(10)

	config.DefaultB = cfg.Section("check").Key("DEFAULT_B").MustFloat64(2)
	config.DefaultK = cfg.Section("check").Key("DEFAULT_K").MustInt(8)
	config.HeapLength = cfg.Section("check").Key("HEAP_LENGTH").MustInt(10)

	config.RedisADDR = cfg.Section("redis").Key("ADDR").MustString("localhost:6379")
	config.RedisPwd = cfg.Section("redis").Key("PWD").MustString("")
	config.DB = cfg.Section("redis").Key("DB").MustInt(0)

	config.MaxConn = cfg.Section("sql").Key("MAX_CONN").MustInt(100)
	config.MaxOpen = cfg.Section("sql").Key("MAX_OPEN").MustInt(10)
	config.SQL = cfg.Section("sql").Key("SQL").MustString("sqllite")
	config.Host = cfg.Section("sql").Key("HOST").MustString("")
	config.User = cfg.Section("sql").Key("USER").MustString("")
	config.Password = cfg.Section("sql").Key("PASSWORD").MustString("")
	config.DBName = cfg.Section("sql").Key("DBNAME").MustString("")

	config.LogPath = cfg.Section("log").Key("PATH").MustString("./log")

	config.Enable = cfg.Section("listen").Key("ENABLE").MustBool(true)
	config.DynamicDuration = cfg.Section("listen").Key("DYNAMIC_DURATION").MustInt64(5)

	config.CacheFilePath = cfg.Section("cache").Key("PATH").MustString("./cache.dat")

	config.RPCPath = cfg.Section("rpc").Key("PATH").MustString("127.0.0.1:8080")
	config.RPCEnable = cfg.Section("rpc").Key("ENABLE").MustBool(false)

	return nil
}
