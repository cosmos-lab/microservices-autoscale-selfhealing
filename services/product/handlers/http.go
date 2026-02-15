package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var Products = []Product{{ID: "1", Name: "Laptop"}}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	for _, p := range Products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
}

func AddProduct(c *gin.Context) {
	var p Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Products = append(Products, p)
	c.JSON(http.StatusCreated, p)
}
