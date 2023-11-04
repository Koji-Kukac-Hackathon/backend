package parkingSpot

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindAllParkingSpotsController(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, FindAllparkingSpotsService())

}

func AddParkingSpotController(ctx *gin.Context) {

	var parkigSpot ParkingSpot

	err := ctx.Bind(&parkigSpot)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	parkingSpot, err := AddParkingSpotService(&parkigSpot)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, parkingSpot)

}
