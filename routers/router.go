package routers

import (
	"fmt"

	"github.com/tsubasa597/ASoulCnkiBackend/check"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
	v1 "github.com/tsubasa597/ASoulCnkiBackend/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	fmt.Println(db.OpenDB())
	check.Init()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apiV1 := router.Group("api/v1")
	{
		apiV1.GET("/check", v1.Check)
		apiV1.GET("/update", v1.Update)
		apiV1.GET("/status", v1.Satus)
	}
	return router
}
