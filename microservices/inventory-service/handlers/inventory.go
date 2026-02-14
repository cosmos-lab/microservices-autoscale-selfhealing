package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// 2️⃣ Async point-to-point
func UpdateStock(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "itemId": "item-123",
        "status": "reserved",
        "message": "Stock updated (async point-to-point)",
    })
}

// 8️⃣ Event streaming / CQRS read
func ReadModel(c *gin.Context) {
    itemId := c.Param("id")
    c.JSON(http.StatusOK, gin.H{
        "itemId": itemId,
        "available": 50,
        "message": "Inventory read model / streaming",
    })
}

