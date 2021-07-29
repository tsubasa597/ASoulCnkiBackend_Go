package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	"github.com/tsubasa597/BILIBILI-HELPER/log"
)

func Update(ctx *gin.Context) {
	go comment.Update(logrus.NewEntry(log.NewLog()))

	ctx.JSON(http.StatusOK, gin.H{
		"code": "ok",
	})
}
