package server

import (
	"fmt"
	"net/http"

	"github.com/tsubasa597/ASoulCnkiBackend/internal/routers"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

var (
	HttpServer *http.Server
)

func NewHttpServer() error {
	eng, err := routers.New()
	if err != nil {
		return err
	}

	HttpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: eng,
	}
	return nil
}
