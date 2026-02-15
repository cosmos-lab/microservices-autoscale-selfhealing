package main

import (
	"inventory-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	go handlers.StartOrderConsumer()

	r := gin.Default()
	r.GET("/inventory/:productId", handlers.GetInventory)
	r.Run(":8080")
}


