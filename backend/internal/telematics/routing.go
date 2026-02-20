package telematics

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/twpayne/go-polyline"
)

// The expected response from OSRM API
type osrmResponse struct {
	Routes []struct {
		Geometry string `json:"geometry"`
	} `json:"routes"`
}

// Wrapper for API call to OSRM routing
func getRoute(start, end []float64) ([]Coordinates, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=polyline",
		start[1], start[0],
		end[1], end[0],
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OSRM request failed: %s", resp.Status)
	}

	var data osrmResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	coords, _, err := polyline.DecodeCoords([]byte(data.Routes[0].Geometry))
	if err != nil {
		return nil, err
	}

	result := make([]Coordinates, 0, len(coords))
	for _, c := range coords {
		result = append(result, Coordinates{
			Lat: c[0],
			Lon: c[1],
		})
	}

	return result, nil
}

func getRouteFromStops(stops []Stop) ([]Coordinates, error) {
	var route []Coordinates

	for i := 0; i < len(stops)-1; i++ {
		from := stops[i].Coords
		to := stops[i+1].Coords

		segment, err := getRoute(
			[]float64{from.Lat, from.Lon},
			[]float64{to.Lat, to.Lon},
		)
		if err != nil {
			return nil, err
		}

		route = append(route, segment...)
	}

	return route, nil
}

// helper methods for stop locating
const earthRadius = 6371000 // meters

func haversine(a, b Coordinates) float64 {
	lat1 := a.Lat * math.Pi / 180
	lon1 := a.Lon * math.Pi / 180
	lat2 := b.Lat * math.Pi / 180
	lon2 := b.Lon * math.Pi / 180

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	sinLat := math.Sin(dLat / 2)
	sinLon := math.Sin(dLon / 2)

	h := sinLat*sinLat +
		math.Cos(lat1)*math.Cos(lat2)*sinLon*sinLon

	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}

func distanceMeters(a, b Coordinates) float64 {
	const R = 6371000 // Earth radius in meters

	lat1 := a.Lat * math.Pi / 180
	lon1 := a.Lon * math.Pi / 180
	lat2 := b.Lat * math.Pi / 180
	lon2 := b.Lon * math.Pi / 180

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	h := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)

	return 2 * R * math.Asin(math.Sqrt(h))
}
