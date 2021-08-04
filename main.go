package main

import (
	"fmt"
	"net/http"

	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/routers"
)

func main() {
	conf.WriteLog()
	comment.InitCache()
	fmt.Println("All Done...")

	eng := routers.Init()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8000),
		Handler: eng,
	}
	s.ListenAndServe()
}
