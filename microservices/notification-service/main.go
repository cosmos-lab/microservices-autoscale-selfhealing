package main

import (
    "github.com/gin-gonic/gin"
    "notification-service/handlers"
)

func main() {
    r := gin.Default()

    // Pub/Sub endpoint
    r.POST("/notification/send", handlers.SendNotification)

    r.Run(":8083")
}

