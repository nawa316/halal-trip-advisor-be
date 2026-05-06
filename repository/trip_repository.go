package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type tripRepository struct {
	db *sql.DB
}

func NewTripRepository(db *sql.DB) domain.TripRepository {
	return &tripRepository{db: db}
}

func (tr *tripRepository) Create(c context.Context, trip *domain.Trip) error {
	_, err := tr.db.ExecContext(c,
		`INSERT INTO trips (id, user_id, name, start_time, end_time, total_distance) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		trip.ID,
		trip.UserID,
		trip.Name,
		trip.StartTime,
		trip.EndTime,
		trip.TotalDistance,
	)
	return err
}

func (tr *tripRepository) FetchByUserID(c context.Context, userID string) ([]domain.Trip, error) {
	rows, err := tr.db.QueryContext(c, `SELECT id, user_id, name, start_time, end_time, total_distance FROM trips WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := make([]domain.Trip, 0)
	for rows.Next() {
		var t domain.Trip
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.StartTime, &t.EndTime, &t.TotalDistance); err != nil {
			return nil, err
		}
		trips = append(trips, t)
	}
	return trips, nil
}

func (tr *tripRepository) GetByID(c context.Context, id string) (domain.Trip, error) {
	var t domain.Trip
	err := tr.db.QueryRowContext(c, `SELECT id, user_id, name, start_time, end_time, total_distance FROM trips WHERE id = $1`, id).
		Scan(&t.ID, &t.UserID, &t.Name, &t.StartTime, &t.EndTime, &t.TotalDistance)
	if err != nil {
		fmt.Printf("Error in Repository GetByID for ID %s: %v\n", id, err)
	}
	return t, err
}

func (tr *tripRepository) Update(c context.Context, trip *domain.Trip) error {
	_, err := tr.db.ExecContext(c,
		`UPDATE trips SET name = $1, start_time = $2, end_time = $3, total_distance = $4, updated_at = NOW() WHERE id = $5`,
		trip.Name,
		trip.StartTime,
		trip.EndTime,
		trip.TotalDistance,
		trip.ID,
	)
	return err
}

func (tr *tripRepository) Delete(c context.Context, id string) error {
	_, err := tr.db.ExecContext(c, `DELETE FROM trips WHERE id = $1`, id)
	return err
}
