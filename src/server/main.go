package server

import (
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
	"zgrabi-mjesto.hr/backend/src/config"
	"zgrabi-mjesto.hr/backend/src/entities/auth"
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

	go func() {
		parkingSpot.AddAllParkingSpotController(out)
	}()

	parkingSpot.Init()
	auth.Init()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Run() {
	databaseTablesInit()
	parkingSpot.Consume()

	r := gin.New()

	r.Use(CORSMiddleware())

	{
		authRoute := r.Group("/auth")
		authRoute.POST("/login", auth.Controller.Login)
		authRoute.POST("/register", auth.Controller.Register)
		authRoute.GET("/user", auth.RequireAuthMiddleware(), auth.Controller.GetUser)
	}

	{
		adminRoute := r.Group("/admin")
		adminRoute.Use(auth.RequireAdminMiddleware())
		adminRoute.GET("/users", auth.Controller.ListUsers)
		adminRoute.POST("/users/:id", auth.Controller.UpdateUser)
		adminRoute.DELETE("/users/:id", auth.Controller.DeleteUser)
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

	r.GET("/parking-spot/filters", parkingSpot.GetParkingSpotsWithFilters)

	r.Run(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port))
}
