package main

import (
	"fmt"

	"zgrabi-mjesto.hr/backend/src/config"
	"zgrabi-mjesto.hr/backend/src/entities/product"
	"zgrabi-mjesto.hr/backend/src/providers/database"
	"zgrabi-mjesto.hr/backend/src/server"
)

func main() {
	config.Config.Init()

	fmt.Printf("Config: %+v\n", config.Config)

	err := database.DatabaseProvider().Register()
	if err != nil {
		panic(err)
	}

	db := database.DatabaseProvider().Client()

	product.Service.Test()

	// Read

	var dbProduct product.Model
	db.First(&dbProduct, 1)                 // find product with integer primary key
	db.First(&dbProduct, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&dbProduct).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&dbProduct).Updates(product.Model{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&dbProduct).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&dbProduct, 1)

	server.Run()
}
