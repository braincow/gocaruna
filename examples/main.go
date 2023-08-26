package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/braincow/gocaruna/caruna"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api, err := caruna.NewCarunaClient()
	if err != nil {
		log.Fatalf("while building api client: %v", err)
	}

	err = api.Login(os.Getenv("CARUNA_USERNAME"), os.Getenv("CARUNA_PASSWORD"))
	if err != nil {
		log.Fatalf("error while doing a login: %v", err)
	}
	defer api.Logout()

	customerIds := api.LoginInfo().User.RepresentedCustomerNumbers

	date := time.Date(2023, time.August, 23, 0, 0, 0, 0, time.Local)

	for _, customerId := range customerIds {
		info, err := api.CustomerInfo(customerId)
		if err != nil {
			log.Fatalf("while querying customer info for %s: %v", customerId, err)
		}
		fmt.Printf("[!!] Customer: %s %s\n", info.ID, info.Name)

		meteringPoints, err := api.MeteringPoints(info.ID)
		if err != nil {
			log.Fatalf("error while querying metering points of the customer: %v", err)
		}
		for _, meteringPoint := range meteringPoints {
			fmt.Printf("[**] Metering point: %s %s\n", meteringPoint.AssetID, meteringPoint.Address.StreetName)

			datas, err := api.ConsumedHours(info.ID, meteringPoint.AssetID, date)
			if err != nil {
				log.Fatalf("error while querying usage data for metering point: %v", err)
			}

			for _, data := range datas {
				fmt.Printf("[..] %s %f kWh\n", data.Timestamp, data.TotalConsumption)
			}
		}
	}
}

// eof
