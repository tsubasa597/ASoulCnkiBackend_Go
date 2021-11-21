package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/request"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/response"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/service/rank"
)

// Rank 作文展
func Rank(ctx *gin.Context) {
	var (
		r request.Rank
	)

	if err := ctx.ShouldBindQuery(&r); err != nil {
		ctx.JSON(http.StatusOK, response.Fail(response.ParamErr))
		return
	}

	ctx.JSON(http.StatusOK, rank.Do(r))
}
