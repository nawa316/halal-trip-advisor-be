package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type placeRepository struct {
	db *sql.DB
}

func NewPlaceRepository(db *sql.DB) domain.PlaceRepository {
	return &placeRepository{db: db}
}

func (pr *placeRepository) Create(c context.Context, place *domain.Place) error {
	_, err := pr.db.ExecContext(c,
		`INSERT INTO places (id, name, latitude, longitude, type, category, rating, open_time, closed_time) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		place.ID,
		place.Name,
		place.Latitude,
		place.Longitude,
		place.Type,
		place.Category,
		place.Rating,
		place.OpenTime,
		place.ClosedTime,
	)
	return err
}

func (pr *placeRepository) Fetch(c context.Context) ([]domain.Place, error) {
	rows, err := pr.db.QueryContext(c, `SELECT id, name, latitude, longitude, type, category, rating, open_time, closed_time FROM places`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	places := make([]domain.Place, 0)
	for rows.Next() {
		var p domain.Place
		if err := rows.Scan(&p.ID, &p.Name, &p.Latitude, &p.Longitude, &p.Type, &p.Category, &p.Rating, &p.OpenTime, &p.ClosedTime); err != nil {
			return nil, err
		}
		places = append(places, p)
	}
	return places, nil
}

func (pr *placeRepository) GetByID(c context.Context, id string) (domain.Place, error) {
	var p domain.Place
	err := pr.db.QueryRowContext(c, 
		`SELECT id, name, latitude, longitude, type, category, rating, open_time, closed_time FROM places WHERE id = $1`, 
		id).Scan(&p.ID, &p.Name, &p.Latitude, &p.Longitude, &p.Type, &p.Category, &p.Rating, &p.OpenTime, &p.ClosedTime)
	if errors.Is(err, sql.ErrNoRows) {
		return p, err
	}
	return p, err
}

func (pr *placeRepository) Update(c context.Context, place *domain.Place) error {
	_, err := pr.db.ExecContext(c,
		`UPDATE places SET name = $1, latitude = $2, longitude = $3, type = $4, category = $5, rating = $6, open_time = $7, closed_time = $8, updated_at = NOW() WHERE id = $9`,
		place.Name,
		place.Latitude,
		place.Longitude,
		place.Type,
		place.Category,
		place.Rating,
		place.OpenTime,
		place.ClosedTime,
		place.ID,
	)
	return err
}

func (pr *placeRepository) Delete(c context.Context, id string) error {
	_, err := pr.db.ExecContext(c, `DELETE FROM places WHERE id = $1`, id)
	return err
}
