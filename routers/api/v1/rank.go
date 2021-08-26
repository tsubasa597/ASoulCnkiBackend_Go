package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
)

func Rank(ctx *gin.Context) {
	var (
		page, size int
		err        error
	)
	if page, err = strconv.Atoi(ctx.Query("page")); err != nil && page < 1 {
		page = 1
	}

	if size, err = strconv.Atoi(ctx.Query("size")); err != nil && (size < 1 || size > 30) {
		size = conf.Size
	}

	ctx.JSON(http.StatusOK, comment.GetInstance().Rank.Do(page, size, ctx.Query("time"), ctx.Query("sort"), ctx.QueryArray("ids")...))
}
