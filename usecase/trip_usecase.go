package usecase

import (
	"context"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type tripUsecase struct {
	tripRepository      domain.TripRepository
	tripRouteRepository domain.TripRouteRepository
	placeRepository     domain.PlaceRepository
	contextTimeout      time.Duration
}

func NewTripUsecase(tripRepository domain.TripRepository, tripRouteRepository domain.TripRouteRepository, placeRepository domain.PlaceRepository, timeout time.Duration) domain.TripUsecase {
	return &tripUsecase{
		tripRepository:      tripRepository,
		tripRouteRepository: tripRouteRepository,
		placeRepository:     placeRepository,
		contextTimeout:      timeout,
	}
}

func (tu *tripUsecase) Create(c context.Context, trip *domain.Trip) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	// 1. Create Trip
	err := tu.tripRepository.Create(ctx, trip)
	if err != nil {
		return err
	}

	// 2. Create Routes
	for _, route := range trip.Routes {
		err := tu.tripRouteRepository.Create(ctx, &route)
		if err != nil {
			// In production, you might want to handle partial failure/rollback
			return err
		}
	}

	return nil
}

func (tu *tripUsecase) GetByID(c context.Context, id string) (domain.Trip, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	trip, err := tu.tripRepository.GetByID(ctx, id)
	if err != nil {
		return domain.Trip{}, err
	}

	routes, err := tu.tripRouteRepository.FetchByTripID(ctx, id)
	if err != nil {
		return trip, nil // Return basic trip if routes fail
	}

	itinerary := make([]domain.TripRouteDetail, 0)
	for _, r := range routes {
		place, _ := tu.placeRepository.GetByID(ctx, r.PlaceID)
		itinerary = append(itinerary, domain.TripRouteDetail{
			TripID:     r.TripID,
			PlaceID:    r.PlaceID,
			OrderIndex: r.OrderIndex,
			Place:      place,
		})
	}
	trip.Itinerary = itinerary

	return trip, nil
}

func (tu *tripUsecase) FetchByUserID(c context.Context, userID string) ([]domain.Trip, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.tripRepository.FetchByUserID(ctx, userID)
}

func (tu *tripUsecase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	// 1. Delete associated routes first (cascading delete if not handled by DB)
	err := tu.tripRouteRepository.DeleteByTripID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Delete trip
	return tu.tripRepository.Delete(ctx, id)
}
