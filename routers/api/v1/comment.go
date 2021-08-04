package v1

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

func Satus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"started": lis.Started(),
		"wait":    atomic.LoadInt32(lis.Wait),
	})
}

func Update(ctx *gin.Context) {
	if lis.Started() {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "fail",
		})
		return
	}

	go lis.Update.All()

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
		lis = comment.NewListen(logrus.NewEntry(logrus.StandardLogger()))

		for _, user := range *db.Get(&db.User{}).(*[]db.User) {
			lis.Add(user)
		}
	})
}
