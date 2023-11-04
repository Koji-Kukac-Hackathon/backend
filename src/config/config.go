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
	Port  int
	Host  string
	DbUrl string
}

var Config = appConfig{}

func (appConfig *appConfig) Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file: ", err)
	}

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

	appConfig.DbUrl = database_url
	appConfig.Host = host
	appConfig.Port = port
}
