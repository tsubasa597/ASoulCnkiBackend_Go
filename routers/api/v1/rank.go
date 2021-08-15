package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Rank(ctx *gin.Context) {
	var (
		t, sort string
	)
	ids := ctx.QueryArray("ids")
	switch ctx.Query("time") {
	case "0":
		t = "0"
	case "1":
		t = fmt.Sprint(time.Now().AddDate(0, 0, -7).Unix())
	case "2":
		t = fmt.Sprint(time.Now().AddDate(0, 0, -3).Unix())
	}
	fmt.Println(t)

	switch ctx.Query("sort") {
	case "1":
		sort = "total_like desc"
	case "2":
		sort = "like desc"
	case "3":
		sort = "num desc"
	}

	data, err := comm.Rank.Do(t, sort, ids...)
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
