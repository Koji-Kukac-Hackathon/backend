package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
	"zgrabi-mjesto.hr/backend/src/entities/product"
	"zgrabi-mjesto.hr/backend/src/providers/database"
)

type appConfig struct {
	Port  int
	Host  string
	DbUrl string
}

const defaultPort = 3000
const defaultHost = "0.0.0.0"

func loadConfig() appConfig {
	envPort, err := strconv.ParseInt(os.Getenv("PORT"), 0, 32)
	if err != nil || envPort == 0 {
		envPort = defaultPort
	}

	envHost := os.Getenv("HOST")
	if envHost == "" {
		envHost = defaultHost
	}

	var port int
	var host string
	var database_url string

	flag.IntVarP(&port, "port", "p", int(envPort), "Set the port on which the server will run")
	flag.StringVarP(&host, "host", "h", envHost, "Set the host to which the server will bind")
	flag.StringVar(&database_url, "database_url", os.Getenv("DATABASE_URL"), "Set the database url")
	flag.Parse()

	return appConfig{
		Port:  port,
		Host:  host,
		DbUrl: database_url,
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file: ", err)
	}

	conf := loadConfig()

	fmt.Printf("Config: %+v\n", conf)

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
}
