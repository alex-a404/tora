package telematics

import "math"

type StopManifest struct {
	Stops []Stop
}

func (sm StopManifest) GetStopFromId(id int32) (*Stop, bool) {
	for _, s := range sm.Stops {
		if s.Id == id {
			return &s, true
		}
	}
	return nil, false
}

func FindNearestStop(
	point Coordinates,
	stops []Stop,
	maxDistance float64, // meters
) (*Stop, bool) {

	bestDist := math.MaxFloat64
	bestID := -1

	for _, stop := range stops {
		d := haversine(point, stop.Coords)

		if d < bestDist {
			bestDist = d
			bestID = int(stop.Id)
		}
	}

	if bestDist > maxDistance {
		return -1, false
	}

	return bestID, true
}

func NewStopManifest_FromYaml(yaml []byte) (*StopManifest, error) {

}
