package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartServer(db *gorm.DB) {
	// ReleaseMode  ช่วยให้เร็วขึ้น แต่จะ debug ไม่ได้
	DB = db
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Load  Controllers
	DemoController(router)
	ProductsController(router)
	CustomerController(router)
	router.Run()

}
