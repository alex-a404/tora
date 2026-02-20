package tracker

import "tora/backend/internal/telematics"

type TrackingSession struct {
	Id       int
	FromStop telematics.Stop
	ToStop   telematics.Stop
	Bus      telematics.Bus
}
