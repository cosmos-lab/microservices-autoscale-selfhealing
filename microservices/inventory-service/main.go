package main

import (
    "github.com/gin-gonic/gin"
    "inventory-service/handlers"
)

func main() {
    r := gin.Default()

    // Async update stock
    r.POST("/inventory/update-stock", handlers.UpdateStock)

    // Event streaming / CQRS read
    r.GET("/inventory/:id/read-model", handlers.ReadModel)

    r.Run(":8082")
}

