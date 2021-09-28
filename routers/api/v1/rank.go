package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
	"github.com/tsubasa597/ASoulCnkiBackend/service/rank"
)

// Rank 作文展
func Rank(ctx *gin.Context) {
	var (
		page, size int
		err        error
	)
	if page, err = strconv.Atoi(ctx.Query("page")); err != nil && page < 1 {
		page = 1
	}

	if size, err = strconv.Atoi(ctx.Query("size")); err != nil && (size < 1 || size > 30) {
		size = setting.Size
	}

	ctx.JSON(http.StatusOK, rank.Do(page, size, ctx.Query("time"), ctx.Query("sort"), ctx.QueryArray("ids")...))
}
