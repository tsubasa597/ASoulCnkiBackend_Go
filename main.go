package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
