package main

import (
	"fmt"
	"net/http"

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
	s.ListenAndServe()
	s.RegisterOnShutdown(func() {

	})
}
