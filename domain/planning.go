package domain

import (
	"context"
)

type PlanningRequest struct {
	StartLat     float64  `json:"start_lat" binding:"required"`
	StartLong    float64  `json:"start_long" binding:"required"`
	StartTime    int64    `json:"start_time" binding:"required"` // Unix timestamp
	EndTime      int64    `json:"end_time" binding:"required"`   // Unix timestamp
	Preferences  []string `json:"preferences"`                    // e.g., "restaurant", "mosque", "tourism"
	MaxPlaces     int      `json:"max_places"`
	Seed          int      `json:"seed"`                           // 0 = optimal, >0 = variations
	ReturnToStart bool     `json:"return_to_start"`
}

type ScheduledPlace struct {
	Place         Place   `json:"place"`
	ArrivalTime   int64   `json:"arrival_time"`
	DepartureTime int64   `json:"departure_time"`
	Distance      float64 `json:"distance_from_previous"` // in km
	ActivityLabel string  `json:"activity_label"`
}

type PlanningResponse struct {
	Itinerary     []ScheduledPlace `json:"itinerary"`
	TotalDistance float64          `json:"total_distance"`
	TotalDuration int64            `json:"total_duration"` // in seconds
	AverageRating float64          `json:"average_rating"`
	Geometry      interface{}      `json:"geometry,omitempty"`
}

type PlanningUsecase interface {
	GenerateRecommendation(c context.Context, req *PlanningRequest) (PlanningResponse, error)
}
