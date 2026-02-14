package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// 1️⃣ Synchronous (Request-Response)
func CreateOrder(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "orderId": "order-123",
        "status":  "created",
        "message": "Order created successfully (sync)",
    })
}

// 4️⃣ Request-Reply over Messaging
func PaymentRequest(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "orderId": "order-123",
        "paymentStatus": "pending",
        "message": "Payment request sent (async request-reply)",
    })
}

// 3️⃣ Publish-Subscribe
func PublishOrderCreated(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "orderId": "order-123",
        "event": "OrderCreated",
        "message": "OrderCreated event published (pub/sub)",
    })
}

// 6️⃣ Saga / Event Sourcing step (dummy)
func SagaStep(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "orderId": "order-123",
        "sagaStep": "payment-processed",
        "message": "Saga step executed",
    })
}

// 7️⃣ CQRS - Read model
func ReadModel(c *gin.Context) {
    orderId := c.Param("id")
    c.JSON(http.StatusOK, gin.H{
        "orderId": orderId,
        "status": "created",
        "readModel": gin.H{
            "total": 100,
            "items": 2,
        },
        "message": "CQRS read model",
    })
}

