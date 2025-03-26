package controller

import (
	"fmt"
	"go-test-grom-by-mikkee/dto"
	model "go-test-grom-by-mikkee/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CustomerController(router *gin.Engine) {
	routes := router.Group("/customers")
	{
		routes.GET("/", getAllCustomers)
		routes.POST("/", createCustomer)
		routes.PUT("/:id", updateCustomer)
		routes.DELETE("/:id", deleteCustomer)
		routes.POST("/login", loginCustomer)
		routes.GET("/profile", getCustomerProfile)
		routes.PUT("/profile/address", updateCustomerAddress)
		routes.PUT("/profile/repassword", Repass)
		routes.GET("/getCarts", getCarts)
	}
}

func getAllCustomers(c *gin.Context) {
	var customers []model.Customer
	if err := DB.Find(&customers).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch customers"})
		return
	}
	c.JSON(200, customers)
}

func createCustomer(c *gin.Context) {
	var customer model.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}
	customer.Password = string(hashedPassword)

	if err := DB.Create(&customer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create customer"})
		return
	}
	c.JSON(201, gin.H{"message": "Customer created successfully"})
}

func loginCustomer(c *gin.Context) {
	var input dto.LoginRequest // Using the LoginRequest DTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var customer model.Customer
	if err := DB.Where("email = ?", input.Email).First(&customer).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// Preparing the response using the LoginResponse DTO
	response := dto.LoginResponse{
		Message: "Login successful",
		Customer: dto.CustomerResponse{
			ID:          customer.CustomerID,
			PhoneNumber: customer.PhoneNumber,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Email:       customer.Email,
			Address:     customer.Address,
		},
	}

	// Send the response back as JSON
	c.JSON(200, response)
}

func Repass(c *gin.Context) {
	// สร้างโครงสร้างรับข้อมูลจาก Body
	var input struct {
		ID          int    `json:"id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	// รับข้อมูลจาก Body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// ดึงข้อมูลลูกค้าตาม id จากฐานข้อมูล
	var customer model.Customer
	if err := DB.First(&customer, input.ID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Customer not found"})
		return
	}

	// ตรวจสอบรหัสผ่านเก่ากับที่บันทึกในระบบ
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(401, gin.H{"error": "Old password is incorrect"})
		return
	}

	// แฮชรหัสผ่านใหม่
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash new password"})
		return
	}

	// อัปเดตรหัสผ่านใหม่
	customer.Password = string(hashedPassword)
	if err := DB.Save(&customer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update password"})
		return
	}
	c.JSON(200, gin.H{"message": "Password updated successfully"})
}

func getCustomerProfile(c *gin.Context) {
	id := c.Query("id")
	var customer model.Customer
	if err := DB.First(&customer, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(200, gin.H{
		"id":      customer.CustomerID,
		"name":    customer.FirstName,
		"email":   customer.Email,
		"address": customer.Address,
	})
}

func updateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer model.Customer
	if err := DB.First(&customer, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Customer not found"})
		return
	}
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	if err := DB.Save(&customer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update customer"})
		return
	}
	c.JSON(200, customer)
}

func updateCustomerAddress(c *gin.Context) {
	id := c.Query("id")
	var customer model.Customer
	if err := DB.First(&customer, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Customer not found"})
		return
	}

	var input struct {
		Address string `json:"address"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	customer.Address = input.Address
	if err := DB.Save(&customer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update address"})
		return
	}
	c.JSON(200, gin.H{"message": "Address updated successfully"})
}

func deleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Delete(&model.Customer{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete customer"})
		return
	}
	c.JSON(200, gin.H{"message": "Customer deleted successfully"})
}

func getCarts(c *gin.Context) {
	// รับ customer_id จาก query parameter
	customerID := c.Query("customer_id")
	if customerID == "" {
		c.JSON(400, gin.H{"error": "Missing customer_id"})
		return
	}

	// แปลง customer_id เป็น integer
	customerIDInt, err := strconv.Atoi(customerID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid customer_id"})
		return
	}

	// ดึงข้อมูลรถเข็นทั้งหมดของลูกค้า พร้อมกับข้อมูลรายการสินค้า
	var carts []model.Cart
	if err := DB.Preload("Items").Where("customer_id = ?", customerIDInt).Find(&carts).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch carts"})
		return
	}

	// ดึงข้อมูลสินค้าทั้งหมดเพื่อลดจำนวนการ query ฐานข้อมูล
	var products []model.Product
	if err := DB.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch products"})
		return
	}

	// สร้าง map เพื่อค้นหาสินค้าตาม ID
	productMap := make(map[int]model.Product)
	for _, product := range products {
		productMap[product.ProductID] = product
	}

	// สร้างโครงสร้างข้อมูลสำหรับ response
	var cartResponses []dto.CartResponseDTO
	for _, cart := range carts {
		var items []dto.CartItemResponseDTO
		for _, item := range cart.Items {
			// ค้นหาสินค้าใน map
			product, exists := productMap[item.ProductID]

			responseItem := dto.CartItemResponseDTO{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}

			if exists {
				responseItem.ProductName = product.ProductName
				responseItem.Description = product.Description
				responseItem.PricePerUnit = product.Price
				// คำนวณราคา
				totalPrice := float64(item.Quantity) * product.Price
				formatted, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", totalPrice), 64)
				responseItem.TotalPrice = formatted
			} else {
				responseItem.ProductName = "Unknown Product (ID: " + strconv.Itoa(item.ProductID) + ")"
				responseItem.Description = "This product is no longer available"
				responseItem.PricePerUnit = 0
				responseItem.TotalPrice = 0
			}

			items = append(items, responseItem)
		}

		cartResponses = append(cartResponses, dto.CartResponseDTO{
			CartID: cart.CartID,
			Name:   cart.CartName,
			Items:  items,
		})
	}

	// ส่งข้อมูลกลับไปยัง client
	c.JSON(200, gin.H{"carts": cartResponses})
}
