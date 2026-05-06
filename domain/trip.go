package domain

import (
	"context"
)

type TripRouteDetail struct {
	TripID     string `json:"trip_id"`
	PlaceID    string `json:"place_id"`
	OrderIndex int64  `json:"order_index"`
	Place      Place  `json:"place"`
}

type Trip struct {
	ID            string            `json:"id"`
	UserID        string            `json:"user_id"`
	Name          string            `json:"name"`
	StartTime     int64             `json:"start_time"`
	EndTime       int64             `json:"end_time"`
	TotalDistance float64           `json:"total_distance"`
	Routes        []TripRoute       `json:"routes,omitempty"`
	Itinerary     []TripRouteDetail `json:"itinerary,omitempty"`
}

type SaveTripRequest struct {
	Name          string   `json:"name" binding:"required"`
	StartTime     int64    `json:"start_time" binding:"required"`
	EndTime       int64    `json:"end_time" binding:"required"`
	TotalDistance float64  `json:"total_distance"`
	PlaceIDs      []string `json:"place_ids" binding:"required"`
}

type TripRepository interface {
	Create(c context.Context, trip *Trip) error
	FetchByUserID(c context.Context, userID string) ([]Trip, error)
	GetByID(c context.Context, id string) (Trip, error)
	Update(c context.Context, trip *Trip) error
	Delete(c context.Context, id string) error
}

type TripUsecase interface {
	Create(c context.Context, trip *Trip) error
	FetchByUserID(c context.Context, userID string) ([]Trip, error)
	GetByID(c context.Context, id string) (Trip, error)
	Delete(c context.Context, id string) error
}
