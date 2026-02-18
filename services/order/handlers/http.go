package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var podName = os.Getenv("POD_NAME")

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

	log.Printf("ORDER_CREATE_REQUEST orderId=%s productId=%s qty=%d pod=%s",
		orderID,
		req.ProductID,
		req.Quantity,
		podName,
	)

	event := OrderEvent{
		OrderID:   orderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Timestamp: time.Now().UnixMilli(),
	}

	payload, _ := json.Marshal(event)

	err := PublishToKafka("orders", orderID, payload)
	if err != nil {
		log.Printf("ORDER_PUBLISH_FAILED orderId=%s err=%v pod=%s", orderID, err, podName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish order event"})
		return
	}

	log.Printf("ORDER_PUBLISHED orderId=%s pod=%s", orderID, podName)

	c.JSON(http.StatusOK, gin.H{
		"orderId":     orderID,
		"status":      "confirmed",
		"productId":   req.ProductID,
		"productName": req.ProductName,
		"quantity":    req.Quantity,
		"message":     "Order created and inventory reserved",
	})
}
