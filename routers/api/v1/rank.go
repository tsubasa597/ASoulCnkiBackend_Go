package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Rank(ctx *gin.Context) {
	var (
		t, sort    string
		page, size int
		err        error
	)
	if page, err = strconv.Atoi(ctx.Query("page")); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"res": nil,
		})
		return
	}

	if size, err = strconv.Atoi(ctx.Query("size")); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"res": nil,
		})
		return
	}

	switch ctx.Query("time") {
	case "1":
		t = fmt.Sprint(time.Now().AddDate(0, 0, -7).Unix())
	case "2":
		t = fmt.Sprint(time.Now().AddDate(0, 0, -3).Unix())
	default:
		t = "0"
	}

	switch ctx.Query("sort") {
	case "1":
		sort = "total_like desc"
	case "2":
		sort = "like desc"
	case "3":
		sort = "num desc"
	}

	data, err := comm.Rank.Do(page, size, t, sort, ctx.QueryArray("ids")...)
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
