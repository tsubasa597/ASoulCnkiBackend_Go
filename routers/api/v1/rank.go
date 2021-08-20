package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
)

func Rank(ctx *gin.Context) {
	var (
		page, size int
		err        error
	)
	if page, err = strconv.Atoi(ctx.Query("page")); err != nil && page < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"res": nil,
		})
		return
	}

	if size, err = strconv.Atoi(ctx.Query("size")); err != nil && (size < 1 || size > 30) {
		ctx.JSON(http.StatusOK, gin.H{
			"res": nil,
		})
		return
	}

	data, err := comment.GetInstance().Rank.Do(page, size, ctx.Query("time"), ctx.Query("sort"), ctx.QueryArray("ids")...)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"res": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"res": data,
	})
}
