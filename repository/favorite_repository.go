package repository

import (
	"context"
	"database/sql"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type favoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) domain.FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (fr *favoriteRepository) Create(c context.Context, favorite *domain.Favorite) error {
	_, err := fr.db.ExecContext(c,
		`INSERT INTO favorites (id, user_id, place_id) VALUES ($1, $2, $3)`,
		favorite.ID,
		favorite.UserID,
		favorite.PlaceID,
	)
	return err
}

func (fr *favoriteRepository) FetchByUserID(c context.Context, userID string) ([]domain.Favorite, error) {
	query := `
		SELECT f.id, f.user_id, f.place_id, p.id, p.name, p.latitude, p.longitude, p.type, p.category, p.rating 
		FROM favorites f
		JOIN places p ON f.place_id = p.id
		WHERE f.user_id = $1
	`
	rows, err := fr.db.QueryContext(c, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	favorites := make([]domain.Favorite, 0)
	for rows.Next() {
		var f domain.Favorite
		var p domain.Place
		if err := rows.Scan(&f.ID, &f.UserID, &f.PlaceID, &p.ID, &p.Name, &p.Latitude, &p.Longitude, &p.Type, &p.Category, &p.Rating); err != nil {
			return nil, err
		}
		f.Place = p
		favorites = append(favorites, f)
	}
	return favorites, nil
}

func (fr *favoriteRepository) Delete(c context.Context, userID string, placeID string) error {
	_, err := fr.db.ExecContext(c, `DELETE FROM favorites WHERE user_id = $1 AND place_id = $2`, userID, placeID)
	return err
}

func (fr *favoriteRepository) IsFavorite(c context.Context, userID string, placeID string) (bool, error) {
	var exists bool
	err := fr.db.QueryRowContext(c, `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND place_id = $2)`, userID, placeID).Scan(&exists)
	return exists, err
}
