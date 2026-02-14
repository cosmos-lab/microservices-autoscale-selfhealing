package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// 3️⃣ Pub/Sub
func SendNotification(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "notificationId": "notif-123",
        "status": "sent",
        "message": "Notification sent (pub/sub)",
    })
}

