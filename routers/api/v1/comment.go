package v1

import (
	"net/http"
	"sync/atomic"
	"time"

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
	comm *comment.Comment
)

func init() {
	db, err := db.New()
	if err != nil {
		panic(err)
	}

	cache, err := cache.New()
	if err != nil {
		panic(err)
	}

	comm = comment.New(*db, *cache, logrus.NewEntry(logrus.StandardLogger()).
		WithField("Time", time.Now().Format("2006/01/02 15:04:05")))

	logrus.WithField("Func", "init").WithField("Time", time.Now().Format("2006/01/02 15:04:05")).
		Info("All Done...")
}
