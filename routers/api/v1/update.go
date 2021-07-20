package v1

import (
	"fmt"

	"github.com/tsubasa597/ASoulCnkiBackend/comments"
	"github.com/tsubasa597/ASoulCnkiBackend/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/log"
)

func Update(ctx *gin.Context) {
	users := db.Get(&db.User{}).(*[]db.User)
	fmt.Println(users)
	for _, v := range *users {
		comments.Update(v, logrus.NewEntry(log.NewLog()))
	}
}
