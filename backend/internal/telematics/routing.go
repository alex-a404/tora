package telematics

import (
	"encoding/json"
	"fmt"
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

func getRouteFromStops(stops [][]float64) ([]Coordinates, error) {
	var route []Coordinates

	for i := 0; i < len(stops)-1; i++ {
		segment, err := getRoute(
			stops[i][1:],   // lat, lon
			stops[i+1][1:], // lat, lon
		)
		if err != nil {
			return nil, err
		}
		route = append(route, segment...)
	}

	return route, nil
}
