package tracker

import "tora/backend/internal/telematics"

// TrackingSession intended to be communicated to client via REST/websockets ...
type TrackingSession struct {
	Id             int
	FromStopName   string
	ToStopName     string
	FromStopCoords telematics.Coordinates
	ToStopCoords   telematics.Coordinates
	Bus            telematics.Bus
}

func (t TrackingSession) Update() {

}
