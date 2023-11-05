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
		"price":                parkingSpot.Price,
	}

	db.Model(ParkingSpot{}).Where("id = ?", parkingSpot.Id).Updates(updateData)

}

func GetAllParkingSpotsWithFilters(zone string, priceMin float32, priceMax float32, latMin float32, latMax float32, longMin float32, longMax float32) ([]ParkingSpot, error) {

	db := database.DatabaseProvider().Client()

	var parkingSpots []ParkingSpot

	if zone != "" {
		// err := db.Table("parking_spots").Where("'parkig_spot_zone' = ?", zone).Where("price >= ? and ? >= price", priceMin, priceMax).Where("latitude >= ? and latitude <= ?", latMin, latMax).Where("longitude >= ? and longitude <= ?", longMin, longMax).Find(&parkingSpots)

		err := db.Model(ParkingSpot{}).Where("parkig_spot_zone = ?", zone).Where("price >= ? and ? >= price", priceMin, priceMax).Find(&parkingSpots)
		if err.Error != nil {
			return []ParkingSpot{}, err.Error
		}
	} else {
		err := db.Model(ParkingSpot{}).Where("price < ?", priceMax).Where("price > ?", priceMin).Find(&parkingSpots)
		if err.Error != nil {
			return []ParkingSpot{}, err.Error
		}
	}

	return parkingSpots, nil

}
