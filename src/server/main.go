package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"zgrabi-mjesto.hr/backend/src/config"
)

func Run() {
	r := gin.New()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port))
}
