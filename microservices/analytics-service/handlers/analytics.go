package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// 3️⃣ Pub/Sub / 8️⃣ Streaming
func ConsumeEvent(c *gin.Context) {
    event := c.Param("event")
    c.JSON(http.StatusOK, gin.H{
        "event": event,
        "status": "consumed",
        "message": "Event consumed by analytics service",
    })
}

