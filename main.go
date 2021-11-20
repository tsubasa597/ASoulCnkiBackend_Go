package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/internal/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/dao"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/server"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/service/listen"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

func main() {
	// 系统日志写入
	log, err := config.NewLogFile("system")
	if err != nil {
		panic(err)
	}

	// 初始化数据库
	if err := dao.Setup(); err != nil {
		panic(err)
	}

	// 初始化缓存
	if err := cache.Setup(); err != nil {
		panic(err)
	}
	defer func() {
		if err := cache.GetCache().Stop(); err != "" {
			log.Error(err)
		}
	}()

	// 初始化自动更新
	listen, err := listen.Setup()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := listen.Stop(); err != nil {
			log.Error(err)
		}
	}()

	// 初始化 http 监听
	s, err := server.NewHttpServer()
	if err != nil {
		panic(err)
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	// 退出后处理
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Error(err)
	}
}
