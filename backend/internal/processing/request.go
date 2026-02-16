package processing

import (
	"math"
	"tora/backend/internal/telematics"
	"tora/backend/internal/tracker"
)

type UserRequest struct {
	FromCoords telematics.Coordinates
	ToCoords   telematics.Coordinates
}

func NewUserRequest(fromCoords, toCoords telematics.Coordinates) UserRequest {
	return UserRequest{
		FromCoords: fromCoords,
		ToCoords:   toCoords,
	}
}

// ProcessReq interacts with python data abstraction layer
func (req UserRequest) ProcessReq(stops []telematics.Stop,
	mgr telematics.Manager,
	trs tracker.TrackingService) UserResponse {

	// determine FROM stop_id and TO stop_id
	fromID, ok := telematics.FindNearestStop(req.FromCoords, stops, 300)
	if !ok {
		return UserResponse{
			false, telematics.Stop{}, telematics.Stop{}, "No available stop", 0,
		}
	}

	toID, ok := telematics.FindNearestStop(req.ToCoords, stops, 300)
	if !ok {
		return UserResponse{
			false, telematics.Stop{}, telematics.Stop{}, "No available stop", 0,
		}
	}

	// TODO find closest bus to FROM
	bus := mgr.GetBuses()[0]

	transfers := bus.Dependencies

	return UserResponse{}
}

type UserResponse struct {
	Ok         bool
	FromStop   telematics.Stop
	ToStop     telematics.Stop
	Message    string
	TrackingID int
}
