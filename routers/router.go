package routers

import (
	"github.com/tsubasa597/ASoulCnkiBackend/comment"
	v1 "github.com/tsubasa597/ASoulCnkiBackend/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.MaxMultipartMemory = 100 << 20

	apiV1 := router.Group("api/v1")
	{
		apiV1.GET("/update", v1.Update)
		apiV1.GET("/check", v1.Check)
		apiV1.GET("/status", v1.Satus)
	}
	comment.InitCache()
	return router
}
