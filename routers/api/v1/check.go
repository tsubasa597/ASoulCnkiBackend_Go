package v1

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/comment/check"
)

func Check(ctx *gin.Context) {
	text := check.ReplaceStr(ctx.PostForm("text"))
	if utf8.RuneCountInString(text) < 8 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "参数错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"res": comment.GetInstance().Compare(text),
	})
}
