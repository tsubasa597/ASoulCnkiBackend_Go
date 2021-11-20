package routers

import (
	"io"

	"github.com/gin-gonic/gin"
	v1 "github.com/tsubasa597/ASoulCnkiBackend/internal/handler/v1"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

func New() (*gin.Engine, error) {
	router := gin.New()

	logger, err := config.NewLogFile("gin")
	if err != nil {
		return nil, err
	}

	// 日志记录
	gin.DisableConsoleColor()
	router.Use(gin.Recovery(), gin.LoggerWithConfig(gin.LoggerConfig{
		Output: io.MultiWriter(logger.Writer()),
	}))

	apiV1 := router.Group("api/v1")
	{
		apiV1.POST("/check", v1.Check)
		apiV1.GET("/rank", v1.Rank)
	}

	return router, nil
}
