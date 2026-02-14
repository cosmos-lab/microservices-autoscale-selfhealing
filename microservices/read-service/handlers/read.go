package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func ReadOrder(c *gin.Context) {
    orderId := c.Param("id")
    c.JSON(http.StatusOK, gin.H{
        "orderId": orderId,
        "status": "created",
        "readModel": gin.H{
            "total": 100,
            "items": 2,
        },
        "message": "Read service - order data (CQRS)",
    })
}

func ReadInventory(c *gin.Context) {
    itemId := c.Param("id")
    c.JSON(http.StatusOK, gin.H{
        "itemId": itemId,
        "available": 50,
        "message": "Read service - inventory data (CQRS)",
    })
}

