package v1

import (
	"net/http"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/check"

	"github.com/gin-gonic/gin"
)

func Check(ctx *gin.Context) {
	text := ctx.Query("text")
	if utf8.RuneCountInString(text) < 8 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "参数错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"res": check.Compare(text),
	})
}
