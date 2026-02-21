package telematics

type Manager struct {
	Buses []Bus
}

func (m *Manager) GetBuses() []Bus {
	return m.Buses
}

func (m *Manager) UpdateAll() {
	for i := range m.Buses {
		m.Buses[i].Update()
	}
}

func NewManager(buses []Bus) *Manager {
	return &Manager{
		Buses: buses,
	}
}
