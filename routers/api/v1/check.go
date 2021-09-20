package v1

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/models/vo"
	"github.com/tsubasa597/ASoulCnkiBackend/service/check"
)

func Check(ctx *gin.Context) {
	text := ctx.PostForm("text")
	if utf8.RuneCountInString(text) < 8 {
		ctx.JSON(http.StatusOK, vo.Fail(vo.TEXTTOSHORT))
		return
	}

	ctx.JSON(http.StatusOK, vo.Sucess(check.Compare(text)))
}
