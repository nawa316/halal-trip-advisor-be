package domain

import (
	"context"
)

type Place struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Type       string  `json:"type"`
	Category   string  `json:"category"`
	Rating     float64 `json:"rating"`
	OpenTime   int64   `json:"open_time"`
	ClosedTime int64   `json:"closed_time"`
}

type PlaceRepository interface {
	Create(c context.Context, place *Place) error
	Fetch(c context.Context) ([]Place, error)
	GetByID(c context.Context, id string) (Place, error)
	Update(c context.Context, place *Place) error
	Delete(c context.Context, id string) error
}
