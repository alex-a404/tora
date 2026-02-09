package main

import "tora/backend/internal/telematics"

var (
	EleftheriaCoords = telematics.Coordinates{Lat: 35.17022784728593, Lon: 33.35889554051766}
	EleftheriaStop   = telematics.Stop{Name: "Eleftheria Sq.", Coords: EleftheriaCoords}

	S1EndCoords = telematics.Coordinates{Lat: 35.13160429484031, Lon: 33.299296813161504}
	S1EndStop   = telematics.Stop{Name: "S1 Initial end", Coords: S1EndCoords}

	S2EndCoords = telematics.Coordinates{Lat: 35.11338633948102, Lon: 33.33255319068168}
	S2EndStop   = telematics.Stop{Name: "S2 Initial end", Coords: S2EndCoords}

	S3EndCoords = telematics.Coordinates{Lat: 35.12984920055877, Lon: 33.36299761491542}
	S3EndStop   = telematics.Stop{Name: "S3 Initial end", Coords: S3EndCoords}

	S4EndCoords = telematics.Coordinates{Lat: 35.14966690544324, Lon: 33.41059674652208}
	S4EndStop   = telematics.Stop{Name: "S4 Initial end", Coords: S3EndCoords}

	S1InitialRoute = []telematics.Coordinates{
		EleftheriaCoords,
		S1EndCoords,
	}

	S2InitialRoute = []telematics.Coordinates{
		EleftheriaCoords,
		S2EndCoords,
	}
	S3InitialRoute = []telematics.Coordinates{
		EleftheriaCoords,
		S3EndCoords,
	}
	S4InitialRoute = []telematics.Coordinates{
		EleftheriaCoords,
		S4EndCoords,
	}
)
