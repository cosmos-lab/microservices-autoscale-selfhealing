package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var products = []Product{
	{ID: "1", Name: "Laptop"},
}

func main() {
	r := gin.Default()

	r.GET("/product/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, p := range products {
			if p.ID == id {
				c.JSON(http.StatusOK, p)
				return
			}
		}
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
