package controller

import (
	model "go-test-grom-by-mikkee/models"

	"github.com/gin-gonic/gin"
)

func CustomerController(router *gin.Engine) {
	routes := router.Group("/customers")
	{
		routes.GET("/", getAllCustomers)

	}
}

func getAllCustomers(c *gin.Context) {
	var customers []model.Customers
	if err := DB.Find(&customers).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(200, customers)
}
