package v1

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/db/vo"
)

func State(ctx *gin.Context) {
	instance := comment.GetInstance()

	ctx.JSON(http.StatusOK, vo.Sucess(vo.Status{
		Enable:  instance.Enable,
		Started: instance.Started(),
		Wait:    atomic.LoadInt32(instance.Wait),
	}))
}

func Update(ctx *gin.Context) {
	instance := comment.GetInstance()
	if instance.Started() {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "fail",
		})
		return
	}

	go instance.Update.All()

	ctx.JSON(http.StatusOK, gin.H{
		"code": "ok",
	})
}
