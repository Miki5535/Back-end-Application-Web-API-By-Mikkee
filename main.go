package main

import (
	"fmt"
	"go-test-grom-by-mikkee/controller"
	model "go-test-grom-by-mikkee/models"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("mysql.dsn"))
	dsn := viper.GetString("mysql.dsn")

	dialactor := mysql.Open(dsn)
	db, err = gorm.Open(dialactor, &gorm.Config{})
	// _, err = gorm.Open(dialactor)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection successful")

	// products := []model.Products{}

	// db.Find(&products)
	// fmt.Println(products)

	model.MigrateModels(db)

	controller.StartServer(db)

}
