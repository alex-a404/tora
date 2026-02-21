package processing

import (
	"tora/backend/internal/telematics"
)

type UserRequest struct {
	FromCoords telematics.Coordinates `json:"from"`
	ToCoords   telematics.Coordinates `json:"to"`
}

type UserResponse struct {
	Ok         bool
	Message    string
	TrackingID string
}

func NewUserRequest(fromCoords, toCoords telematics.Coordinates) UserRequest {
	return UserRequest{
		FromCoords: fromCoords,
		ToCoords:   toCoords,
	}
}
