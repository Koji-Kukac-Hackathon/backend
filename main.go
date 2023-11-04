package main

import (
	"fmt"

	"zgrabi-mjesto.hr/backend/src/config"
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

	server.Run()
}
