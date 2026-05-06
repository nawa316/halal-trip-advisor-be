package repository

import (
	"context"
	"database/sql"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type tripRouteRepository struct {
	db *sql.DB
}

func NewTripRouteRepository(db *sql.DB) domain.TripRouteRepository {
	return &tripRouteRepository{db: db}
}

func (trr *tripRouteRepository) Create(c context.Context, tripRoute *domain.TripRoute) error {
	_, err := trr.db.ExecContext(c,
		`INSERT INTO trip_routes (id, trip_id, place_id, order_index) 
		 VALUES ($1, $2, $3, $4)`,
		tripRoute.ID,
		tripRoute.TripID,
		tripRoute.PlaceID,
		tripRoute.OrderIndex,
	)
	return err
}

func (trr *tripRouteRepository) FetchByTripID(c context.Context, tripID string) ([]domain.TripRoute, error) {
	rows, err := trr.db.QueryContext(c, `SELECT id, trip_id, place_id, order_index FROM trip_routes WHERE trip_id = $1 ORDER BY order_index ASC`, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routes := make([]domain.TripRoute, 0)
	for rows.Next() {
		var tr domain.TripRoute
		if err := rows.Scan(&tr.ID, &tr.TripID, &tr.PlaceID, &tr.OrderIndex); err != nil {
			return nil, err
		}
		routes = append(routes, tr)
	}
	return routes, nil
}

func (trr *tripRouteRepository) DeleteByTripID(c context.Context, tripID string) error {
	_, err := trr.db.ExecContext(c, `DELETE FROM trip_routes WHERE trip_id = $1`, tripID)
	return err
}
