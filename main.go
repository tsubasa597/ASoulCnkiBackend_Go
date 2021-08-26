package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/routers"
)

func main() {
	conf.WriteLog()

	eng := routers.Init()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: eng,
	}
	go s.ListenAndServe()

	s.RegisterOnShutdown(func() {
		comment.GetInstance().Stop()
	})

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	s.Shutdown(ctx)
}
