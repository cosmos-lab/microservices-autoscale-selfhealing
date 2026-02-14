package main

import (
    "github.com/gin-gonic/gin"
    "analytics-service/handlers"
)

func main() {
    r := gin.Default()

    r.GET("/analytics/event/:event", handlers.ConsumeEvent)

    r.Run(":8084")
}

