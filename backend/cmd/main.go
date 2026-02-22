package main

import (
	"log"
	"time"
	"tora/backend/internal/processing"
	"tora/backend/internal/telematics"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	s processing.Service
)

func routeSetup() {
	// load stop manifest
	manifest, err := telematics.NewStopManifestFromFile("../data/urban_stops.json")
	if err != nil {
		logger.Fatalf(err.Error())
	}
	s.StopManifest = manifest

	EleftheriaStop, _ := s.StopManifest.GetStopFromId(581)
	S1EndStop, _ := s.StopManifest.GetStopFromId(248)
	S2EndStop, _ := s.StopManifest.GetStopFromId(297)
	S3EndStop, _ := s.StopManifest.GetStopFromId(377)
	S4EndStop, _ := s.StopManifest.GetStopFromId(1216)

	// setup buses
	busS1, err := telematics.NewBus(
		"S1",
		telematics.RouteDependency{
			From:       *EleftheriaStop,
			To:         *S1EndStop,
			InProgress: true,
		},
		EleftheriaStop.Coords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS2, err := telematics.NewBus(
		"S2",
		telematics.RouteDependency{
			From:       *EleftheriaStop,
			To:         *S2EndStop,
			InProgress: true,
		},
		EleftheriaStop.Coords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS3, err := telematics.NewBus(
		"S3",
		telematics.RouteDependency{
			From:       *EleftheriaStop,
			To:         *S3EndStop,
			InProgress: true,
		},
		EleftheriaStop.Coords,
	)
	if err != nil {
		log.Fatal(err)
	}

	busS4, err := telematics.NewBus(
		"S4",
		telematics.RouteDependency{
			From:       *EleftheriaStop,
			To:         *S4EndStop,
			InProgress: true,
		},
		EleftheriaStop.Coords,
	)
	if err != nil {
		log.Fatal(err)
	}

	s.TelematicsMgr = telematics.NewManager([]telematics.Bus{*busS1, *busS2, *busS3, *busS4})
	logger.Infof("Inital stops loaded, buses initialised.")
}

func main() {

	routeSetup()
	go updateWorker()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/get_buses", func(c *gin.Context) {
		buses := s.TelematicsMgr.GetBuses()
		c.JSON(200, buses)
	})

	r.GET("/get_session", func(c *gin.Context) {
		if trackingId, ok := c.Params.Get("session_id"); ok {
			c.JSON(200, s.TrackingService.GetResponse(trackingId))
		} else {
			c.JSON(404, "server error session_id not found")
		}
	})

	r.POST("/request_ride", func(c *gin.Context) {
		var rideReq processing.UserRequest

		if err := c.ShouldBindJSON(&rideReq); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		resp := s.ProcessReq(rideReq)

		if resp.Ok {
			c.JSON(200, resp.TrackingID)
		} else {
			c.JSON(414, resp.Message)
		}
	})

	err := r.Run(":8000")
	if err != nil {
		return
	}

}

func updateWorker() {
	for {
		s.TelematicsMgr.UpdateAll()
		time.Sleep(time.Second)
	}
}
