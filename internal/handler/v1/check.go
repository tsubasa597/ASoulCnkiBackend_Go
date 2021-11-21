package v1

import (
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/request"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/response"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/service/check"
)

// Check 查重
func Check(ctx *gin.Context) {
	var c request.Check
	if err := ctx.ShouldBind(&c); err != nil {
		ctx.JSON(http.StatusOK, response.Fail(response.ParamErr))
		return
	}

	if utf8.RuneCountInString(c.Text) < 8 {
		ctx.JSON(http.StatusOK, response.Fail(response.TextShort))
		return
	}

	ctx.JSON(http.StatusOK, response.Sucess(check.Compare(c.Text)))
}
