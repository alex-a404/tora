package telematics

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Stop struct {
	Name   string `json:"name"`
	Coords Coordinates
	Id     int `json:"id"`
}

// RouteDependency a rule that stop From should come before stop To in an itinerary
type RouteDependency struct {
	From       Stop
	To         Stop
	InProgress bool
}

type Bus struct {
	Name         string        `json:"name"`
	Itinerary    []Stop        `json:"itinerary"`
	Pos          Coordinates   `json:"pos"`
	Route        []Coordinates `json:"route"`
	Dependencies []RouteDependency
}

func NewBus(name string, initialTransfers RouteDependency, pos Coordinates) (*Bus, error) {

	b := &Bus{
		Name:      name,
		Itinerary: []Stop{initialTransfers.From, initialTransfers.To},
		Pos:       pos,
	}

	// calculate initial route
	if err := b.CalculateRoute(); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Bus) Redirect(from int, to int) error {
	err := b.CalculateRoute()
	if err != nil {
		return err
	}
	return nil
}

func (b *Bus) CalculateRoute() error {
	route, err := getRouteFromStops(b.Itinerary)
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
	// move bus
	b.Pos = b.Route[0]
	b.Route = b.Route[1:]

	const threshold = 5.0
	newTransfers := make([]RouteDependency, 0, len(b.Dependencies))
	for _, t := range b.Dependencies {
		// Check if we are near pickup stop
		if !t.InProgress &&
			distanceMeters(b.Pos, t.From.Coords) <= threshold {

			t.InProgress = true
		}

		// Check if we are near dropoff stop AND transfer active
		if t.InProgress &&
			distanceMeters(b.Pos, t.To.Coords) <= threshold {

			// Drop completed, do not add to newTransfers
			continue
		}
		newTransfers = append(newTransfers, t)
	}
	b.Dependencies = newTransfers
}
