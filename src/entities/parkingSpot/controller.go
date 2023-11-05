package parkingSpot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func FindAllParkingSpotsController(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, FindAllparkingSpotsService())

}

type ParkingSpotTmp struct {
	Id                 string  `json:"id"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	ParkigSpotZone     string  `json:"parkingSpotZone"`
	Occupied           bool    `json:"occupied"`
	OccupiedTimesStamp string  `json:"occupiedTimestamp"`
}

func AddAllParkingSpotController(data []byte) {
	var parkingSpots []ParkingSpotTmp
	err := json.Unmarshal(data, &parkingSpots)
	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}

	for _, parkingSpot := range parkingSpots {

		fmt.Println(parkingSpot.OccupiedTimesStamp)
		time, err := time.Parse("2006-01-02T15:04:05.999999999", parkingSpot.OccupiedTimesStamp)

		if err != nil {
			fmt.Println("error", err)
		}

		parkingSpotCurrent := ParkingSpot{
			Id:                 parkingSpot.Id,
			Latitude:           parkingSpot.Latitude,
			Longitude:          parkingSpot.Longitude,
			ParkigSpotZone:     parkingSpot.ParkigSpotZone,
			Occupied:           parkingSpot.Occupied,
			OccupiedTimesStamp: time,
		}

		UpdateParkingSpotService(&parkingSpotCurrent)
	}
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
