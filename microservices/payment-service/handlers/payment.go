package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// 1️⃣ Synchronous (Request-Response)
func ProcessPayment(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "paymentId": "pay-123",
        "status": "approved",
        "message": "Payment approved (sync)",
    })
}

// 4️⃣ Request-Reply over Messaging
func PaymentRequestReply(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "paymentId": "pay-123",
        "status": "processed",
        "message": "Payment processed (async request-reply)",
    })
}

