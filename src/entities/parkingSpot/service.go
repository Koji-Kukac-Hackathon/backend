package parkingSpot

import (
	"errors"
	"fmt"

	"zgrabi-mjesto.hr/backend/src/providers/database"
)

func FindAllparkingSpotsService() []ParkingSpot {
	db := database.DatabaseProvider().Client()

	var parkingSpot []ParkingSpot

	db.Find(&parkingSpot)
	return parkingSpot
}

// if parsing is not successful returns an error and empty ParkingSpot struct
func AddParkingSpotService(parkingSpot *ParkingSpot) (ParkingSpot, error) {
	db := database.DatabaseProvider().Client()

	var err error

	fmt.Println(*parkingSpot)

	if err != nil {
		return ParkingSpot{}, errors.New("unable to parse timestamp")
	}

	db.Create(&parkingSpot)

	return *parkingSpot, nil
}
