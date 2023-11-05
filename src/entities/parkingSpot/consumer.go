package parkingSpot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

type ParkingSpotScraped struct {
	Id         string `json:"Id"`
	IsOccupied bool   `json:"IsOccupied"`
	Time       string `json:"Time"`
}

func Consume() {

	hub, err := eventhub.NewHubFromConnectionString(os.Getenv("EVENT_HUB_CONNECTION_STRING"), eventhub.HubWithSenderMaxRetryCount(5))

	if err != nil {
		fmt.Println("\n\n\n\nhub init failed")
		panic(err)
	}

	exit := make(chan struct{})

	var buffer [90000]ParkingSpotScraped
	buffpointer := 0

	ticker := time.NewTicker(2 * time.Second)

	// Creating channel using make
	tickerChan := make(chan bool)
	go func() {
		for {
			select {
			case <-tickerChan:
				return
			// interval task
			case <-ticker.C:
				fmt.Println("exec")
				for i := 0; i < buffpointer; i++ {

					parkingSpot, err := GetParkingSpotService(buffer[i].Id)

					if err != nil {
						fmt.Println(err)
					}

					fmt.Println("--------", parkingSpot)
					parkingSpot.Occupied = buffer[i].IsOccupied
					fmt.Println(parkingSpot)
					parkingSpot.LastDataReceived = time.Now()
					if parkingSpot.Occupied == true {
						parkingSpot.OccupiedTimesStamp = time.Now()
					}
					fmt.Println(parkingSpot)
					UpdateParkingSpotService(&parkingSpot)
				}
				buffpointer = 0
			}
		}
	}()

	handler := func(ctx context.Context, event *eventhub.Event) error {
		text := string(event.Data)
		if text == "exit\n" {
			fmt.Println("Oh snap!! Someone told me to exit!")
			exit <- *new(struct{})
		} else {
			var data ParkingSpotScraped
			err := json.Unmarshal(event.Data, &data)
			if err != nil {
				fmt.Println("error", err)
			}
			fmt.Println("got data")
			fmt.Println(string(event.Data))
			buffer[buffpointer] = data
			buffpointer++
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_, err = hub.Receive(ctx, "0", handler, eventhub.ReceiveWithLatestOffset())
	cancel()

	if err != nil {

		fmt.Println("\n\n\n\nhub receive failed")
		panic(err)
	}

}
