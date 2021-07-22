package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
)

func Satus(ctx *gin.Context) {
	started, wait := comment.Status()
	ctx.JSON(http.StatusOK, gin.H{
		"started": started,
		"wait":    wait,
	})
}
