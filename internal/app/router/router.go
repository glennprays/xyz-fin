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
	transactionHandler *handler.TransactionHandler,
) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	apiV1 := router.Group(BasePath)

	apiV1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	consumerGroup := apiV1.Group("/consumers")
	{
		consumerGroup.POST("/login", consumerHandler.Login)
		consumerGroup.GET("/:nik", authMiddleware.Authenticate(), consumerHandler.GetByNIK)
	}

	apiV1.POST("/transactions", authMiddleware.Authenticate(), transactionHandler.CreateTransaction)

	return router
}
