package main

import (
	"tora/backend/internal/telematics"

	"github.com/gin-gonic/gin"
)

var (
	m telematics.Manager
)

func routeSetup() {
	// setup buses
	busS1 := telematics.Bus{
		Name: "S1",
		Pos:  EleftheriaCoords,
		Stops: []telematics.Stop{
			EleftheriaStop,
			S1EndStop,
		},
		Route: S1InitialRoute,
	}

	busS2 := telematics.Bus{
		Name: "S2",
		Pos:  EleftheriaCoords,
		Stops: []telematics.Stop{
			EleftheriaStop,
			S2EndStop,
		},
		Route: S2InitialRoute,
	}

	busS3 := telematics.Bus{
		Name: "S3",
		Pos:  EleftheriaCoords,
		Stops: []telematics.Stop{
			EleftheriaStop,
			S3EndStop,
		},
		Route: S3InitialRoute,
	}

	busS4 := telematics.Bus{
		Name: "S4",
		Pos:  EleftheriaCoords,
		Stops: []telematics.Stop{
			EleftheriaStop,
			S1EndStop,
		},
		Route: S4InitialRoute,
	}

	m = telematics.Manager{
		Buses: []telematics.Bus{busS1, busS2, busS3, busS4},
	}
}

func main() {

	routeSetup()

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
