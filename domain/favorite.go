package domain

import (
	"context"
)

type Favorite struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	PlaceID string `json:"place_id"`
	Place   Place  `json:"place,omitempty"`
}

type FavoriteRepository interface {
	Create(c context.Context, favorite *Favorite) error
	FetchByUserID(c context.Context, userID string) ([]Favorite, error)
	Delete(c context.Context, userID string, placeID string) error
	IsFavorite(c context.Context, userID string, placeID string) (bool, error)
}

type FavoriteUsecase interface {
	Create(c context.Context, favorite *Favorite) error
	FetchByUserID(c context.Context, userID string) ([]Favorite, error)
	Delete(c context.Context, userID string, placeID string) error
	IsFavorite(c context.Context, userID string, placeID string) (bool, error)
}
