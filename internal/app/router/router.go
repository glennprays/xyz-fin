package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glennprays/xyz-fin/internal/app/handler"
	"github.com/glennprays/xyz-fin/internal/app/middleware"
)

const BasePath = "/api/v1"

func SetupRouter(
	authMiddleware *middleware.AuthMiddleware,
	consumerHandler *handler.ConsumerHandler,
) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	consumerGroup := router.Group(BasePath + "/consumers")
	{
		consumerGroup.POST("/login", consumerHandler.Login)
	}

	return router
}
