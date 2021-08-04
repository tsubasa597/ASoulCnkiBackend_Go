package v1

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

func Satus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"started": lis.Update.Started || lis.Started(),
		"wait":    atomic.LoadInt32(lis.Wait),
	})
}

func Update(ctx *gin.Context) {
	if lis.Update.Started || lis.Started() {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "fail",
		})
		return
	}

	go lis.Update.Do(*db.Get(&db.User{}).(*[]db.User))

	ctx.JSON(http.StatusOK, gin.H{
		"code": "ok",
	})
}

var (
	once sync.Once
	lis  *comment.Listen
)

func init() {
	once.Do(func() {
		conf.WriteLog()
		lis = comment.NewListen(logrus.NewEntry(logrus.StandardLogger()))
	})
}
