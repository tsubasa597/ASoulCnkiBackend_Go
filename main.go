package main

import (
	"fmt"
	"net/http"

	"github.com/tsubasa597/ASoulCnkiBackend/routers"
)

func main() {
	eng := routers.Init()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8000),
		Handler: eng,
	}
	s.ListenAndServe()
}
