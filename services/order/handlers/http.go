package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderRequest struct {
	ProductID   string `json:"productId" binding:"required"`
	ProductName string `json:"productName" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
}

type OrderEvent struct {
	OrderID   string `json:"orderId"`
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	Timestamp int64  `json:"timestamp"`
}

type InventoryEvent struct {
	OrderID   string `json:"orderId"`
	ProductID string `json:"productId"`
	NewStock  int    `json:"newStock,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

func CreateOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID := "order-" + time.Now().Format("20060102150405")

	event := OrderEvent{
		OrderID:   orderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Timestamp: time.Now().UnixMilli(),
	}

	payload, _ := json.Marshal(event)

	err := PublishToKafka("orders", orderID, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish order event"})
		return
	}

	result, err := WaitForInventoryReply(orderID, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "inventory service timeout"})
		return
	}

	if result.Reason != "" {
		c.JSON(http.StatusConflict, gin.H{
			"orderId": orderID,
			"status":  "failed",
			"reason":  result.Reason,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderId":     orderID,
		"status":      "confirmed",
		"productId":   req.ProductID,
		"productName": req.ProductName,
		"quantity":    req.Quantity,
		"newStock":    result.NewStock,
		"message":     "Order created and inventory reserved",
	})
}
