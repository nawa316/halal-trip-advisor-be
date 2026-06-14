package geoutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type MapboxRouteLeg struct {
	Distance float64 `json:"distance"` // in meters
	Duration float64 `json:"duration"` // in seconds
}

type MapboxRoute struct {
	Distance float64          `json:"distance"`
	Duration float64          `json:"duration"`
	Legs     []MapboxRouteLeg `json:"legs"`
	Geometry interface{}      `json:"geometry,omitempty"`
}

type MapboxResponse struct {
	Code   string        `json:"code"`
	Routes []MapboxRoute `json:"routes"`
}

type Coordinate struct {
	Lat  float64
	Long float64
}

func GetMapboxRoute(coords []Coordinate) (*MapboxRoute, error) {
	token := viper.GetString("MAPBOX_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("MAPBOX_API_KEY is not set")
	}

	if len(coords) < 2 {
		return nil, fmt.Errorf("at least 2 coordinates required")
	}

	// Mapbox uses Longitude,Latitude order!
	coordStrings := make([]string, len(coords))
	for i, c := range coords {
		coordStrings[i] = fmt.Sprintf("%f,%f", c.Long, c.Lat)
	}

	coordPath := strings.Join(coordStrings, ";")
	url := fmt.Sprintf("https://api.mapbox.com/directions/v5/mapbox/driving-traffic/%s?access_token=%s&geometries=geojson", coordPath, token)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mapbox api returned status %d", resp.StatusCode)
	}

	var mapboxResp MapboxResponse
	if err := json.NewDecoder(resp.Body).Decode(&mapboxResp); err != nil {
		return nil, err
	}

	if mapboxResp.Code != "Ok" || len(mapboxResp.Routes) == 0 {
		return nil, fmt.Errorf("mapbox api returned no routes, code: %s", mapboxResp.Code)
	}

	return &mapboxResp.Routes[0], nil
}
