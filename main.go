package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/logging"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/model"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
	"github.com/tsubasa597/ASoulCnkiBackend/routers"
	"github.com/tsubasa597/ASoulCnkiBackend/service/listen"
)

func main() {
	logging.Write()

	model.Setup()
	cache.Setup()
	listen := listen.New(logrus.NewEntry(logrus.StandardLogger()))

	eng := routers.Init()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Port),
		Handler: eng,
	}
	go s.ListenAndServe()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	s.Shutdown(ctx)

	cache.GetCache().Stop()
	listen.Stop()
}
