package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderRequest struct {
	ProductID   string `json:"productId" binding:"required"`
	ProductName string `json:"productName" binding:"required"`
}

func CreateOrder(c *gin.Context) {
	var req OrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderId":     "order-123",
		"status":      "created",
		"productName": req.ProductName,
		"productId":   req.ProductID,
		"message":     "Order created successfully (sync)",
	})
}
