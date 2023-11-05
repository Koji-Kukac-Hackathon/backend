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

func GetParkingSpotService(id string) (ParkingSpot, error) {
	db := database.DatabaseProvider().Client()

	var parkingSpot ParkingSpot

	retval := db.Table("parking_spots").Where("id = ?", id).Order("occupied_times_stamp DESC").First(&parkingSpot)

	if retval.Error != nil {
		return ParkingSpot{}, errors.New("no such parking spot")
	}

	return parkingSpot, nil
}

func UpdateParkingSpotService(parkingSpot *ParkingSpot) {
	db := database.DatabaseProvider().Client()

	var updateData map[string]interface{}
	updateData = map[string]interface{}{
		"occupied":             parkingSpot.Occupied,
		"occupied_times_stamp": parkingSpot.OccupiedTimesStamp,
		"last_data_received":   parkingSpot.LastDataReceived,
	}

	db.Model(ParkingSpot{}).Where("id = ?", parkingSpot.Id).Updates(updateData)

}
func UpdateParkingSpotServiceWithGeoData(parkingSpot *ParkingSpot) {
	db := database.DatabaseProvider().Client()

	var updateData map[string]interface{}
	updateData = map[string]interface{}{
		"longitude":            parkingSpot.Longitude,
		"latitude":             parkingSpot.Latitude,
		"occupied":             parkingSpot.Occupied,
		"occupied_times_stamp": parkingSpot.OccupiedTimesStamp,
	}

	db.Model(ParkingSpot{}).Where("id = ?", parkingSpot.Id).Updates(updateData)

}

func GetAllParkingSpotsWithSameZone(zone string) ([]ParkingSpot, error) {

	db := database.DatabaseProvider().Client()

	var parkingSpots []ParkingSpot

	err := db.Model(ParkingSpot{}).Where("parkig_spot_zone = ?", zone).Find(&parkingSpots)

	if err.Error != nil {
		return []ParkingSpot{}, err.Error
	}

	return parkingSpots, nil

}

func GetAllParkingSpotsWithOccupation(occupation bool) ([]ParkingSpot, error) {
	db := database.DatabaseProvider().Client()

	var parkingSpots []ParkingSpot

	err := db.Model(ParkingSpot{}).Where("occupied =?", occupation).Find(&parkingSpots)

	if err.Error != nil {
		return []ParkingSpot{}, err.Error
	}

	return parkingSpots, nil
}
