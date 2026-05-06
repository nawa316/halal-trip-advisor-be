package domain

import (
	"context"
)

type TripRoute struct {
	ID         string `json:"id"`
	TripID     string `json:"trip_id"`
	PlaceID    string `json:"place_id"`
	OrderIndex int64  `json:"order_index"`
}

type TripRouteRepository interface {
	Create(c context.Context, tripRoute *TripRoute) error
	FetchByTripID(c context.Context, tripID string) ([]TripRoute, error)
	DeleteByTripID(c context.Context, tripID string) error
}
