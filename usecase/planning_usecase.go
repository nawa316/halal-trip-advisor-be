package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/internal/geoutil"
)

type planningUsecase struct {
	placeRepository domain.PlaceRepository
	contextTimeout  time.Duration
}

func NewPlanningUsecase(placeRepository domain.PlaceRepository, timeout time.Duration) domain.PlanningUsecase {
	return &planningUsecase{
		placeRepository: placeRepository,
		contextTimeout:  timeout,
	}
}

func (pu *planningUsecase) GenerateRecommendation(c context.Context, req *domain.PlanningRequest) (domain.PlanningResponse, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	allPlaces, err := pu.placeRepository.Fetch(ctx)
	if err != nil {
		return domain.PlanningResponse{}, err
	}

	// 1. Filter places by preference and basic availability
	filteredPlaces := []domain.Place{}
	prefMap := make(map[string]bool)
	for _, p := range req.Preferences {
		prefMap[strings.ToLower(p)] = true
	}

	for _, p := range allPlaces {
		if len(req.Preferences) > 0 {
			match := false
			// Check Type
			if _, ok := prefMap[strings.ToLower(p.Type)]; ok {
				match = true
			}
			// Check Category
			if !match {
				if _, ok := prefMap[strings.ToLower(p.Category)]; ok {
					match = true
				}
			}
			// Special handling for keywords like "halal"
			if !match && prefMap["halal"] {
				if strings.Contains(strings.ToLower(p.Type), "halal") || strings.Contains(strings.ToLower(p.Category), "halal") {
					match = true
				}
			}
			// Special handling for "tourism"
			if !match && prefMap["tourism"] {
				if strings.Contains(strings.ToLower(p.Category), "tourist") {
					match = true
				}
			}

			if match {
				filteredPlaces = append(filteredPlaces, p)
			}
		} else {
			filteredPlaces = append(filteredPlaces, p)
		}
	}

	fmt.Printf("Total places: %d, Filtered places: %d\n", len(allPlaces), len(filteredPlaces))

	// 2. Greedy approach based on distance
	itinerary := []domain.ScheduledPlace{}
	currentTime := req.StartTime
	currentLat := req.StartLat
	currentLong := req.StartLong
	totalDistance := 0.0

	// Track visited places
	visited := make(map[string]bool)

	for len(itinerary) < req.MaxPlaces && currentTime < req.EndTime {
		var nextPlace *domain.Place
		minDist := -1.0
		
		for i := range filteredPlaces {
			p := &filteredPlaces[i]
			if visited[p.ID] {
				continue
			}

			dist := geoutil.Haversine(currentLat, currentLong, p.Latitude, p.Longitude)
			
			// Simple check: can we reach it?
			// Assume travel speed 30km/h for buffer
			travelTimeSeconds := int64((dist / 30.0) * 3600)
			arrivalTime := currentTime + travelTimeSeconds
			
			if arrivalTime > req.EndTime {
				continue
			}

			// Scoring: Lower score is better.
			// Formula: (Normalized Distance * 0.7) + (Normalized Inverse Rating * 0.3)
			// Normalized Distance: dist / 50.0 (capped at 1.0)
			// Normalized Inverse Rating: (5.0 - p.Rating) / 5.0
			
			normDist := dist / 50.0
			if normDist > 1.0 {
				normDist = 1.0
			}
			normRating := (5.0 - p.Rating) / 5.0
			
			score := (normDist * 0.6) + (normRating * 0.4)

			if nextPlace == nil || score < minDist { // reusing minDist variable for score check
				minDist = score
				nextPlace = p
			}
		}

		if nextPlace == nil {
			break 
		}

		// Recalculate actual distance for the selected place
		actualDist := geoutil.Haversine(currentLat, currentLong, nextPlace.Latitude, nextPlace.Longitude)

		// Calculate arrival and departure
		travelTimeSeconds := int64((actualDist / 30.0) * 3600)
		arrivalTime := currentTime + travelTimeSeconds
		departureTime := arrivalTime + 3600 // Spend 1 hour

		itinerary = append(itinerary, domain.ScheduledPlace{
			Place:                 *nextPlace,
			ArrivalTime:           arrivalTime,
			DepartureTime:         departureTime,
			Distance:              actualDist,
		})

		visited[nextPlace.ID] = true
		currentTime = departureTime
		currentLat = nextPlace.Latitude
		currentLong = nextPlace.Longitude
		totalDistance += actualDist
	}

	return domain.PlanningResponse{
		Itinerary:     itinerary,
		TotalDistance: totalDistance,
		TotalDuration: currentTime - req.StartTime,
	}, nil
}
