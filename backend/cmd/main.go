package main

import (
	"fmt"
	"log"
	"time"
	"tora/backend/internal/telematics"

	"github.com/gin-gonic/gin"
)

var (
	m telematics.Manager
)

func routeSetup() {
	// setup buses
	busS1, err := telematics.NewBus(
		"S1",
		[]telematics.Stop{
			EleftheriaStop,
			S1EndStop,
		},
		EleftheriaCoords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS2, err := telematics.NewBus(
		"S2",
		[]telematics.Stop{
			EleftheriaStop,
			S2EndStop,
		},
		EleftheriaCoords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS3, err := telematics.NewBus(
		"S3",
		[]telematics.Stop{
			EleftheriaStop,
			S3EndStop,
		},
		EleftheriaCoords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS4, err := telematics.NewBus(
		"S4",
		[]telematics.Stop{
			EleftheriaStop,
			S4EndStop,
		},
		EleftheriaCoords,
	)
	if err != nil {
		log.Fatal(err)
	}

	m = telematics.Manager{
		Buses: []telematics.Bus{*busS1, *busS2, *busS3, *busS4},
	}
}

func main() {

	routeSetup()
	go updateWorker()

	r := gin.Default()
	r.GET("/get_buses", func(c *gin.Context) {
		buses := m.GetBuses()
		c.JSON(200, buses)
	})
	err := r.Run(":8080")
	if err != nil {
		return
	}

}

func requestTransfer(c *gin.Context) {}
func updateWorker() {
	for {
		fmt.Println("updateWorker")
		m.UpdateAll()
		time.Sleep(time.Second)
	}
}
