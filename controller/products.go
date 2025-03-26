package controller

import (
	"errors"
	"go-test-grom-by-mikkee/dto"
	model "go-test-grom-by-mikkee/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductsController(router *gin.Engine) {
	routes := router.Group("/products")
	{
		routes.GET("/", getAllProducts)
		routes.POST("/", createProduct)
		routes.PUT("/:id", updateProduct)
		routes.DELETE("/:id", deleteProduct)
		routes.POST("/addProductToCart", addProductToCart)
		routes.GET("/search", searchProducts)
	}
}

func getAllProducts(c *gin.Context) {
	var products []model.Product
	if err := DB.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch products"})
		return
	}

	// สร้าง DTO สำหรับแต่ละ Product
	var productDTOs []dto.ProductDTO
	for _, product := range products {
		productDTO := dto.ProductDTO{
			ProductID:     product.ProductID,
			ProductName:   product.ProductName,
			Description:   product.Description,
			Price:         product.Price,
			StockQuantity: product.StockQuantity,
		}
		productDTOs = append(productDTOs, productDTO)
	}

	c.JSON(200, productDTOs)
}

func createProduct(c *gin.Context) {
	var product model.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := DB.Create(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create product"})
		return
	}

	// สร้าง DTO
	productDTO := dto.ProductDTO{
		ProductID:     product.ProductID,
		ProductName:   product.ProductName,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
	}

	c.JSON(201, productDTO)
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var product model.Product

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
	var product model.Product

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

func addProductToCart(c *gin.Context) {
	var input struct {
		CustomerID int    `json:"customer_id"`
		CartName   string `json:"cart_name"`
		ProductID  int    `json:"product_id"`
		Quantity   int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// ตรวจสอบว่าสินค้ามีอยู่ในฐานข้อมูลหรือไม่
	var product model.Product
	if err := DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	// ค้นหารถเข็นตามชื่อและ customer_id
	var cart model.Cart
	if err := DB.Where("customer_id = ? AND cart_name = ?", input.CustomerID, input.CartName).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// สร้างรถเข็นใหม่หากไม่พบ
			cart = model.Cart{
				CustomerID: input.CustomerID,
				CartName:   input.CartName,
			}
			if err := DB.Create(&cart).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			c.JSON(500, gin.H{"error": "Failed to fetch cart"})
			return
		}
	}

	// ตรวจสอบว่าสินค้าอยู่ในรถเข็นแล้วหรือไม่
	var cartItem model.CartItem
	if err := DB.Where("cart_id = ? AND product_id = ?", cart.CartID, input.ProductID).First(&cartItem).Error; err == nil {
		// หากมีสินค้าอยู่แล้ว ให้เพิ่มจำนวน
		cartItem.Quantity += input.Quantity
		if err := DB.Save(&cartItem).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update cart item"})
			return
		}
	} else {
		// หากไม่มีสินค้า ให้เพิ่มรายการใหม่
		cartItem = model.CartItem{
			CartID:    cart.CartID,
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
		}
		if err := DB.Create(&cartItem).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to add product to cart"})
			return
		}
	}

	// ลดจำนวน stock_quantity ของสินค้าที่ถูกเพิ่มในรถเข็น
	product.StockQuantity -= input.Quantity
	if err := DB.Save(&product).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update product stock"})
		return
	}

	c.JSON(200, gin.H{"message": "Product added to cart successfully"})
}

func searchProducts(c *gin.Context) {
	// ดึง query parameters จาก URL
	keyword := c.Query("keyword")
	minPrice := c.Query("minPrice")
	maxPrice := c.Query("maxPrice")

	// สร้าง query เริ่มต้น
	query := DB.Model(&model.Product{})

	// ถ้ามี keyword ให้ค้นหาใน ProductName และ Description
	if keyword != "" {
		query = query.Where("product_name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// ถ้ามี minPrice ให้เพิ่มเงื่อนไขราคาขั้นต่ำ
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}

	// ถ้ามี maxPrice ให้เพิ่มเงื่อนไขราคาขั้นสูง
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	var products []model.Product
	if err := query.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to search products"})
		return
	}

	// แปลงเป็น DTO
	var productDTOs []dto.ProductDTO
	for _, product := range products {
		productDTO := dto.ProductDTO{
			ProductID:     product.ProductID,
			ProductName:   product.ProductName,
			Description:   product.Description,
			Price:         product.Price,
			StockQuantity: product.StockQuantity,
		}
		productDTOs = append(productDTOs, productDTO)
	}

	c.JSON(200, productDTOs)
}
