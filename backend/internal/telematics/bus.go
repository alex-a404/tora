package telematics

type Coordinates struct {
	Lat float64
	Lon float64
}

type Stop struct {
	Name   string
	Coords Coordinates
}

type Bus struct {
	Name  string        `json:"name"`
	Stops []Stop        `json:"stops"`
	Pos   Coordinates   `json:"pos"`
	Route []Coordinates `json:"route"`
}

func NewBus(name string, stops []Stop, pos Coordinates) (*Bus, error) {
	b := &Bus{
		Name:  name,
		Stops: stops,
		Pos:   pos,
	}

	// calculate initial route
	if err := b.CalculateRoute(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bus) CalculateRoute() error {
	route, err := getRouteFromStops(b.Stops)
	if err != nil {
		return err
	}
	b.Route = route
	return nil
}

func (b *Bus) Update() {
	if len(b.Route) == 0 {
		return
	}
	b.Pos = b.Route[0]
	b.Route = b.Route[1:]
}
