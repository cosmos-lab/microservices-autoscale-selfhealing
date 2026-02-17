package main

import (
	"log"
	"order-service/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	podName := os.Getenv("POD_NAME")
	log.Printf("ORDER_SERVICE_STARTING pod=%s", podName)

	handlers.StartInventoryReplyConsumer()

	router := gin.Default()

	router.POST("/order", handlers.CreateOrder)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
	log.Println("ORDER_SERVICE_LISTENING port=" + port)
}
