package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var podName = os.Getenv("POD_NAME")

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var products = []Product{
	{ID: "1", Name: "Laptop"},
}

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("/product/:id", func(c *gin.Context) {

		id := c.Param("id")

		log.Printf("PRODUCT_FETCH_REQUEST productId=%s pod=%s", id, podName)

		for _, p := range products {
			if p.ID == id {
				log.Printf("PRODUCT_FOUND productId=%s pod=%s", id, podName)
				c.JSON(http.StatusOK, p)
				return
			}
		}

		log.Printf("PRODUCT_NOT_FOUND productId=%s pod=%s", id, podName)
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	})

	r.POST("/product", func(c *gin.Context) {
		var newProd Product
		if err := c.ShouldBindJSON(&newProd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products = append(products, newProd)
		c.JSON(http.StatusCreated, newProd)
	})

	r.Run(":8080")
}
