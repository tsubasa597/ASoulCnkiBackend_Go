package routers

import (
	v1 "github.com/tsubasa597/ASoulCnkiBackend/routers/api/v1"

	"github.com/gin-gonic/gin"
)

// Init 初始化 gin
func Init() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apiV1 := router.Group("api/v1")
	{
		apiV1.POST("/check", v1.Check)
		apiV1.GET("/rank", v1.Rank)
	}

	return router
}
