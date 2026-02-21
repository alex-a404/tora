package telematics

import (
	"encoding/json"
	"math"
	"os"
)

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
	var bestStop *Stop

	for _, stop := range stops {
		d := haversine(point, stop.Coords)

		if d < bestDist {
			bestDist = d
			bestStop = &stop
		}
	}

	if bestDist > maxDistance {
		return nil, false
	}

	return bestStop, true
}

func NewStopManifestFromJson(data []byte) (*StopManifest, error) {
	var manifest StopManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func NewStopManifestFromFile(path string) (*StopManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewStopManifestFromJson(data)
}
