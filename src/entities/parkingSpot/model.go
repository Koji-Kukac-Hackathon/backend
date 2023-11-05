package parkingSpot

import (
	"net/http"
	"time"

	"zgrabi-mjesto.hr/backend/src/providers/database"
)

var UnauthorizedParkingSpot = ParkingSpotError{
	Type:     "unauthorized",
	Title:    "Unauthorized",
	Status:   http.StatusUnauthorized,
	Detail:   "You are not authorized to access this parking spot",
	Instance: "ParkingSpot error",
}

var NotFoundParkingSpot = ParkingSpotError{
	Type:     "unauthorized",
	Title:    "Unauthorized",
	Status:   http.StatusNotFound,
	Detail:   "This page does not exist",
	Instance: "ParkingSpot error",
}

type ParkingSpot struct {
	Id                 string    `json:"id"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`
	ParkigSpotZone     string    `json:"parkingSpotZone"`
	Occupied           bool      `json:"occupied"`
	OccupiedTimesStamp time.Time `json:"occupiedTimestamp"`
	LastDataReceived   time.Time `json:"lastDataReceived"`
	Price              float32   `json:"price"`
}

type ParkingSpotError struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func Init() {
	db := database.DatabaseProvider().Client()

	db.AutoMigrate(&ParkingSpot{})
}
