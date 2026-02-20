package main

import (
	"fmt"
	"log"
	"time"
	"tora/backend/internal/telematics"
	"tora/backend/internal/tracker"

	"github.com/gin-gonic/gin"
)

var (
	mgr telematics.Manager
	sm  telematics.StopManifest
	trs tracker.TrackingService
)

func routeSetup() {
	// setup buses
	busS1, err := telematics.NewBus(
		"S1",
		telematics.RouteDependency{
			From:       EleftheriaStop,
			To:         S1EndStop,
			InProgress: true,
		},
		EleftheriaCoords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS2, err := telematics.NewBus(
		"S2",
		telematics.RouteDependency{
			From:       EleftheriaStop,
			To:         S2EndStop,
			InProgress: true,
		},
		EleftheriaCoords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS3, err := telematics.NewBus(
		"S3",
		telematics.RouteDependency{
			From:       EleftheriaStop,
			To:         S3EndStop,
			InProgress: true,
		},
		EleftheriaCoords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS4, err := telematics.NewBus(
		"S4",
		telematics.RouteDependency{
			From:       EleftheriaStop,
			To:         S4EndStop,
			InProgress: true,
		},
		EleftheriaCoords,
	)
	if err != nil {
		log.Fatal(err)
	}

	mgr = telematics.Manager{
		Buses: []telematics.Bus{*busS1, *busS2, *busS3, *busS4},
	}
}

func main() {

	routeSetup()
	go updateWorker()

	r := gin.Default()
	r.GET("/get_buses", func(c *gin.Context) {
		buses := mgr.GetBuses()
		c.JSON(200, buses)
	})
	err := r.Run(":8000")
	if err != nil {
		return
	}

}

func updateWorker() {
	for {
		fmt.Println("updateWorker")
		mgr.UpdateAll()
		trs.UpdateAll()
		time.Sleep(time.Second)
	}
}
