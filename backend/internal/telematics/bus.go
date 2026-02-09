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
	Name   string
	Coords Coordinates
}

type Bus struct {
	Name  string        `json:"name"`
	Stops []Stop        `json:"stops"`
	Pos   Coordinates   `json:"pos"`
	Route []Coordinates `json:"route"`
}

func NewBus(name string, pos string) *Bus {
	// parse initial coords of bus
	parts := strings.Split(pos, ",")
	lat, _ := strconv.ParseFloat(parts[0], 64)
	lon, _ := strconv.ParseFloat(parts[1], 64)

	return &Bus{
		Name: name,
		Pos: Coordinates{
			Lat: lat,
			Lon: lon,
		},
	}
}
