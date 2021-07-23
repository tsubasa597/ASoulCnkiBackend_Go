package v1

import (
	"net/http"

	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/ASoulCnkiBackend/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/log"
)

func Update(ctx *gin.Context) {
	users := db.Get(&db.User{}).(*[]db.User)
	for _, v := range *users {
		go comment.Update(v, logrus.NewEntry(log.NewLog()))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "ok",
	})
}
