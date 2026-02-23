package tracker

import "tora/backend/internal/telematics"

type TrackingSession struct {
	Id       string
	FromStop telematics.Stop
	ToStop   telematics.Stop
	Bus      *telematics.Bus
}
