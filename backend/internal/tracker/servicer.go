package tracker

import "tora/backend/internal/telematics"

type TrackingService struct {
	ActiveSessions map[string]TrackingSession
}

func (t TrackingService) RegisterSession(session TrackingSession) {
	t.ActiveSessions[session.Id] = session
}

func (t TrackingService) UnregisterSession(session TrackingSession) {
	delete(t.ActiveSessions, session.Id)
}

func (t TrackingService) GetSession(sessionId string) TrackingSession {
	return t.ActiveSessions[sessionId]
}

// TrackingSessionResponse intended to be communicated to client via REST/websockets ...
type TrackingSessionResponse struct {
	BusId          string                 `json:"bus_id"`
	FromStopCoords telematics.Coordinates `json:"from_stop_coords"`
	ToStopCoords   telematics.Coordinates `json:"to_stop_coords"`
	FromStopName   string                 `json:"from_stop_name"`
	ToStopName     string                 `json:"to_stop_name"`
	BusCoords      telematics.Coordinates `json:"bus_coords"`
	ETA            int                    `json:"eta"`
}

func (t TrackingService) GetResponse(sessionId string) TrackingSessionResponse {
	tsession := t.GetSession(sessionId)
	// todo implement ETA
	return TrackingSessionResponse{
		BusId:          tsession.Bus.Name,
		FromStopCoords: tsession.FromStop.Coords,
		ToStopCoords:   tsession.ToStop.Coords,
		FromStopName:   tsession.FromStop.Name,
		ToStopName:     tsession.ToStop.Name,
		BusCoords:      tsession.Bus.Pos,
		ETA:            0,
	}
}
