package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"zgrabi-mjesto.hr/backend/src/config"
	"zgrabi-mjesto.hr/backend/src/entities/auth"
	"zgrabi-mjesto.hr/backend/src/entities/parkingSpot"
)

func databaseTablesInit() {
	parkingSpot.Init()
	auth.Init()
}

func Run() {
	databaseTablesInit()

	r := gin.New()

	{
		authRoute := r.Group("/auth")
		authRoute.POST("/login", auth.Controller.Login)
		authRoute.POST("/register", auth.Controller.Register)
		authRoute.GET("/user", auth.RequireAuthMiddleware(), auth.Controller.GetUser)
	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ping-authed", auth.RequireAuthMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ping-admin", auth.RequireAdminMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/parking-spot", parkingSpot.FindAllParkingSpotsController)

	r.POST("/parking-spot", parkingSpot.AddParkingSpotController)

	r.Run(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port))
}
