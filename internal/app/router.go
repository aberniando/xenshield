package app

import (
	"github.com/aberniando/xenshield/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(handler *gin.Engine, handlers *Handlers) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	transactions := handler.Group("/transactions")
	transactions.Use(middleware.APIKeyAuth())
	transactions.GET("/:ip_address", handlers.TransactionHandler.GetIPAddressStolenCardHistory)
	transactions.POST("", handlers.TransactionHandler.Insert)
}
