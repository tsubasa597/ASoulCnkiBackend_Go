package v1

import (
	"github.com/tsubasa597/ASoulCnkiBackend/comments"
	"github.com/tsubasa597/ASoulCnkiBackend/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/log"
)

func Update(ctx *gin.Context) {
	users := db.Get(&db.User{}).(*[]db.User)
	for _, v := range *users {
		go comments.Update(v, logrus.NewEntry(log.NewLog()))
	}
}
