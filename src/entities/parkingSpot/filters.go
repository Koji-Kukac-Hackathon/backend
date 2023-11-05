package parkingSpot

import (
	"fmt"
	"math"
)

type Filters struct {
	ParkingSpotZone string  `form:"parkingSpotZone"`
	PriceMax        float32 `form:"priceMax"`
	PriceMin        float32 `form:"priceMin"`
	Longitude       float64 `form:"longitude"`
	Latitude        float64 `form:"latitude"`
	Radius          float64 `form:"radius"`
}

func FindAllParkingSpotsWithFilters(filters *Filters) ([]ParkingSpot, error) {

	var parkingSpots []ParkingSpot

	fmt.Println(filters)

	latMin := filters.Latitude - (filters.Radius / 110574)
	latMax := filters.Latitude + (filters.Radius / 110574)

	longMin := filters.Longitude + (filters.Radius / (111320 * math.Cos(filters.Latitude)))
	longMax := filters.Longitude - (filters.Radius / (111320 * math.Cos(filters.Latitude)))

	parkingSpots, err := GetAllParkingSpotsWithFilters(filters.ParkingSpotZone, filters.PriceMin, filters.PriceMax, float32(latMin), float32(latMax), float32(longMin), float32(longMax))
	if err != nil {

		return []ParkingSpot{}, err
	}

	return parkingSpots, nil
}
