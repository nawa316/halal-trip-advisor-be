package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	token := os.Getenv("MAPBOX_API_KEY")
	if token == "" {
		fmt.Println("No MAPBOX_API_KEY")
		return
	}
	url := fmt.Sprintf("https://api.mapbox.com/directions/v5/mapbox/driving-traffic/115.168,-8.746;115.178,-8.756;115.168,-8.746?access_token=%s&geometries=geojson&steps=true", token)
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	routes := data["routes"].([]interface{})
	legs := routes[0].(map[string]interface{})["legs"].([]interface{})
	fmt.Printf("Legs count: %d\n", len(legs))
	steps := legs[0].(map[string]interface{})["steps"].([]interface{})
	geom := steps[0].(map[string]interface{})["geometry"]
	fmt.Printf("Step 0 geometry type: %T\n", geom)
	if m, ok := geom.(map[string]interface{}); ok {
		fmt.Printf("Geometry type: %s\n", m["type"])
	}
}
