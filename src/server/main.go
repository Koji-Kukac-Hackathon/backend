package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"zgrabi-mjesto.hr/backend/src/config"
	"zgrabi-mjesto.hr/backend/src/entities/parkingSpot"
)

func databaseTablesInit() {
	parkingSpot.Init()
}

func Run() {

	databaseTablesInit()

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
