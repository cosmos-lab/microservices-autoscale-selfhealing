package main

import (
    "github.com/gin-gonic/gin"
    "order-service/handlers"
)

func main() {
    r := gin.Default()

    // Synchronous (Request-Response)
    r.POST("/order/create", handlers.CreateOrder)

    // Async / Request-Reply over messaging (dummy)
    r.POST("/order/payment-request", handlers.PaymentRequest)

    // Pub/Sub simulation (dummy)
    r.POST("/order/publish-event", handlers.PublishOrderCreated)

    // Saga / Event Sourcing (dummy)
    r.POST("/order/saga-step", handlers.SagaStep)

    // CQRS
    r.GET("/order/:id/read-model", handlers.ReadModel)

    r.Run(":8080") // Listen on port 8080
}

