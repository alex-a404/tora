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

type Transfer struct {
	From       Stop
	To         Stop
	InProgress bool
}

type Bus struct {
	Name      string        `json:"name"`
	Stops     []Stop        `json:"stops"`
	Pos       Coordinates   `json:"pos"`
	Route     []Coordinates `json:"route"`
	Transfers []Transfer
}

func NewBus(name string, initialTransfers Transfer, pos Coordinates) (*Bus, error) {

	b := &Bus{
		Name:  name,
		Stops: []Stop{initialTransfers.From, initialTransfers.To},
		Pos:   pos,
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
	// move bus
	b.Pos = b.Route[0]
	b.Route = b.Route[1:]

	const threshold = 5.0
	newTransfers := make([]Transfer, 0, len(b.Transfers))
	for _, t := range b.Transfers {
		// Check if we are near pickup stop
		if !t.InProgress &&
			distanceMeters(b.Pos, t.From.Coords) <= threshold {

			t.InProgress = true
		}

		// Check if we are near dropoff stop AND transfer active
		if t.InProgress &&
			distanceMeters(b.Pos, t.To.Coords) <= threshold {

			// Drop completed → do NOT add to newTransfers
			continue
		}
		newTransfers = append(newTransfers, t)
	}
	b.Transfers = newTransfers
}
