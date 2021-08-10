package v1

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

func Satus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"started": comm.Started(),
		"wait":    atomic.LoadInt32(comm.Wait),
	})
}

func Update(ctx *gin.Context) {
	if comm.Started() {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "fail",
		})
		return
	}

	go comm.Update.All()

	ctx.JSON(http.StatusOK, gin.H{
		"code": "ok",
	})
}

var (
	once sync.Once
	comm *comment.Comment
)

func init() {
	once.Do(func() {
		db, err := db.New()
		if err != nil {
			panic(err)
		}

		cache, err := cache.NewComment()
		if err != nil {
			panic(err)
		}

		comm = comment.New(*db, cache, logrus.NewEntry(logrus.StandardLogger()))

		// for _, user := range *db.Get(&entry.User{}).(*[]entry.User) {
		// 	comm.Add(user)
		// }

		fmt.Println("All Done...")
	})
}
