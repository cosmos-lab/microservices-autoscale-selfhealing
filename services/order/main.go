package main

import (
	"order-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	handlers.StartInventoryReplyConsumer()

	r := gin.Default()
	r.POST("/order/create", handlers.CreateOrder)
	r.Run(":8080")
}


