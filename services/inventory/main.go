package main

import (
	"inventory-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	go handlers.StartOrderConsumer()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("/inventory/:productId", handlers.GetInventory)
	r.Run(":8080")
}
