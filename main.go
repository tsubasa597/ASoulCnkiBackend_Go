package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/tsubasa597/ASoulCnkiBackend/routers"
)

func main() {
	eng := routers.Init()
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8000),
		Handler: eng,
	}
	s.ListenAndServe()
	// f, err := os.Create("trace.pprof")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer f.Close()

	// if err := trace.Start(f); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer trace.Stop()

	// db.OpenDB()
	// check.Init()
	// users := db.Get(&db.User{}).(*[]db.User)
	// fmt.Println(users)
	// for _, v := range *users {
	// 	comments.Update(v, logrus.NewEntry(log.NewLog()))
	// }
}
