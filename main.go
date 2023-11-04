package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"zgrabi-mjesto.hr/backend/model"
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

	db, err := gorm.Open(mysql.Open(conf.DbUrl), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.Product{})

	// Create
	db.Create(&model.Product{Code: "D42", Price: 100})

	// Read
	var product model.Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(model.Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product, 1)
}
