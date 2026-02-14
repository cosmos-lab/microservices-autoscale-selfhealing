package main

import (
    "github.com/gin-gonic/gin"
    "payment-service/handlers"
)

func main() {
    r := gin.Default()

    // Synchronous payment approval
    r.POST("/payment/process", handlers.ProcessPayment)

    // Async / Request-Reply messaging
    r.POST("/payment/request-reply", handlers.PaymentRequestReply)

    r.Run(":8081")
}

