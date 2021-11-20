package server

import (
	"fmt"
	"net/http"

	"github.com/tsubasa597/ASoulCnkiBackend/internal/routers"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

func NewHttpServer() (*http.Server, error) {
	eng, err := routers.New()
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: eng,
	}, nil
}
