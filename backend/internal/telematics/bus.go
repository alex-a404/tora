package telematics

import (
	"strconv"
	"strings"
)

type Coordinates struct {
	Lat float64
	Lon float64
}

type Stop struct {
	name   string
	coords Coordinates
}

type Bus struct {
	name  string
	stops []Stop
	pos   Coordinates
	route []Coordinates
}

func NewBus(name string, pos string) *Bus {
	// parse initial coords of bus
	parts := strings.Split(pos, ",")
	lat, _ := strconv.ParseFloat(parts[0], 64)
	lon, _ := strconv.ParseFloat(parts[1], 64)

	return &Bus{
		name: name,
		pos: Coordinates{
			Lat: lat,
			Lon: lon,
		},
	}
}
