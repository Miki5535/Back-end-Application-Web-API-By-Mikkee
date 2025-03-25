package controller

import (
	model "go-test-grom-by-mikkee/models"

	"github.com/gin-gonic/gin"
)

func ProductsController(router *gin.Engine) {
	routes := router.Group("/products")
	{
		routes.GET("/", getAllProducts)
		routes.POST("/", createProduct)
		routes.PUT("/:id", updateProduct)
		routes.DELETE("/:id", deleteProduct)
	}
}

func getAllProducts(c *gin.Context) {
	var products []model.Products
	if err := DB.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(200, products)
}

func createProduct(c *gin.Context) {
	var product model.Products

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := DB.Create(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(201, product)
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var product model.Products

	if err := DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := DB.Save(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(200, product)
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product model.Products

	if err := DB.First(&product, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	if err := DB.Delete(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}
