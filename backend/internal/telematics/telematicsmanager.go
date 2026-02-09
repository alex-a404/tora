package telematics

type Manager struct {
	Buses []Bus
}

func (m *Manager) GetBuses() []Bus {
	return m.Buses
}
