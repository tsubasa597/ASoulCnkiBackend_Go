package main

import (
	"fmt"
	"net/http"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/routers"
)

func main() {
	conf.WriteLog()

	// db, err := db.New()
	// if err != nil {
	// 	panic(err)
	// }

	// cache, err := cache.NewComment()
	// if err != nil {
	// 	panic(err)
	// }

	// comm := comment.New(*db, cache, logrus.NewEntry(logrus.StandardLogger()))

	eng := routers.Init()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: eng,
	}
	s.ListenAndServe()
	s.RegisterOnShutdown(func() {

	})
}
