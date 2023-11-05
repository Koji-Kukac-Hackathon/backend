package server

import (
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
	"zgrabi-mjesto.hr/backend/src/config"
	"zgrabi-mjesto.hr/backend/src/entities/parkingSpot"
)

func databaseTablesInit() {

	curl := exec.Command("curl", "https://hackathon.kojikukac.com/api/ParkingSpot/getAll", "-X", "GET", "-H", "accept: application/json", "-H", "Api-Key: b7e43abf-190c-4b66-bb4d-909659e125db")
	out, err := curl.Output()
	if err != nil {
		fmt.Println("erorr", err)
		return
	}
	fmt.Println(string(out))
	parkingSpot.AddAllParkingSpotController(out)

	parkingSpot.Init()
}

func Run() {

	databaseTablesInit()
	parkingSpot.Consume()

	r := gin.New()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/parking-spot", parkingSpot.FindAllParkingSpotsController)

	r.POST("/parking-spot", parkingSpot.AddParkingSpotController)

	r.Run(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port))
}
