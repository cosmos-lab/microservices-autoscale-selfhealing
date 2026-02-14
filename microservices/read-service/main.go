package main

import (
    "github.com/gin-gonic/gin"
    "read-service/handlers"
)

func main() {
    r := gin.Default()

    r.GET("/read/orders/:id", handlers.ReadOrder)
    r.GET("/read/inventory/:id", handlers.ReadInventory)

    r.Run(":8085")
}

