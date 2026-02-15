package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type InventoryItem struct {
	ProductID string `json:"productId"`
	Stock     int    `json:"stock"`
}

var (
	mu        sync.Mutex
	inventory = map[string]int{
		"1": 999999,
		"2": 999999,
	}
)

func GetInventory(c *gin.Context) {
	productID := c.Param("productId")
	mu.Lock()
	stock, ok := inventory[productID]
	mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, InventoryItem{
		ProductID: productID,
		Stock:     stock,
	})
}


