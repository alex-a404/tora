package tracker

type TrackingService struct {
	ActiveSessions []TrackingSession
}

func (t TrackingService) RegisterSession(session TrackingSession) {
	t.ActiveSessions = append(t.ActiveSessions, session)
}

func (t TrackingService) UnregisterSession(session TrackingSession) {}

func (t TrackingService) UpdateAll() {

}
