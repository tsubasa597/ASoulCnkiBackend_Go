package v1

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
)

// 字符串中带有 % 会解析错误
func Check(ctx *gin.Context) {
	text := comment.ReplaceStr(ctx.Query("text"))

	if utf8.RuneCountInString(text) < 8 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "参数错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"res": comment.Compare(text),
	})
}
