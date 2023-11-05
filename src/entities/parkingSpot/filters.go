package parkingSpot

import "fmt"

type Filters struct {
	ParkingSpotZone string  `form:"parkingSpotZone"`
	Occupied        bool    `form:"occupied"`
	PriceMax        float32 `form:"priceMax"`
	PriceMin        float32 `form:"priceMin"`
	Longitude       float64 `form:"longitude"`
	Latitude        float64 `form:"latitude"`
}

func FindAllParkingSpotsWithFilters(filters *Filters) ([]ParkingSpot, error) {

	var parkingSpots []ParkingSpot
	var parkingSpotsTmp []ParkingSpot

	fmt.Println(filters)
	if filters.ParkingSpotZone != "" {
		parkingSpotsTmp, err := GetAllParkingSpotsWithSameZone(filters.ParkingSpotZone)
		if err != nil {
			return []ParkingSpot{}, err
		}
		parkingSpots = append(parkingSpots, parkingSpotsTmp...)
	}

	parkingSpotsTmp, err := GetAllParkingSpotsWithOccupation(filters.Occupied)
	if err != nil {
		return []ParkingSpot{}, err
	}

	parkingSpots = append(parkingSpots, parkingSpotsTmp...)

	return parkingSpots, nil
}
