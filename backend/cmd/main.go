package main

import (
	"log"
	"time"
	"tora/backend/internal/processing"
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

	r.GET("/get_session", func(c *gin.Context) {
		if trackingId, ok := c.Params.Get("session_id"); ok {
			c.JSON(200, trs.GetResponse(trackingId))
		} else {
			c.JSON(404, "server error session_id not found")
		}
	})

	r.POST("/request_ride", func(c *gin.Context) {
		if fromStop, ok := c.Params.Get("from_stop"); ok {

		}

		if toStop, ok := c.Params.Get("to_stop"); ok {

		}
		req := processing.UserRequest{
			FromCoords: fromStop,
			ToCoords:   toStop,
		}
		resp := req.ProcessReq(mgr, sm, trs)

		if resp.Ok {
			c.JSON(200, resp.TrackingID)
		} else {
			c.JSON(414, "An error occured")
		}

	})

	err := r.Run(":8000")
	if err != nil {
		return
	}

}

func updateWorker() {
	for {
		mgr.UpdateAll()
		time.Sleep(time.Second)
	}
}
