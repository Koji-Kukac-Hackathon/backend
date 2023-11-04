package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
)

const defaultPort = 3000
const defaultHost = "0.0.0.0"

type appConfig struct {
	Port      int
	Host      string
	DbUrl     string
	ApiSecret string
}

var Config = appConfig{}

func (appConfig *appConfig) Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file: ", err)
	}

	{
		envPort, err := strconv.ParseInt(os.Getenv("PORT"), 0, 32)
		if err != nil || envPort == 0 {
			envPort = defaultPort
		}
		port := int(envPort)
		flag.IntVarP(&port, "port", "p", int(envPort), "Set the port on which the server will run")
		appConfig.Port = port
	}
	{
		host := os.Getenv("HOST")
		if host == "" {
			host = defaultHost
		}
		flag.StringVarP(&host, "host", "h", host, "Set the host to which the server will bind")
		appConfig.Host = host
	}
	{
		var database_url string
		flag.StringVar(&database_url, "database-url", os.Getenv("DATABASE_URL"), "Set the database url")
		appConfig.DbUrl = database_url
	}
	{
		var apiSecret string
		flag.StringVar(&apiSecret, "api-secret", os.Getenv("API_SECRET"), "Set the api secret")
		appConfig.ApiSecret = apiSecret
	}

	flag.Parse()
}
