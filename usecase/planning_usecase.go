package usecase

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
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

// MealWindow defines a time range for a specific meal.
type MealWindow struct {
	Label    string // "Sarapan 🍳", "Makan Siang 🍽️", "Makan Malam 🌙"
	StartHour int
	StartMin  int
	EndHour   int
	EndMin    int
}

var mealWindows = []MealWindow{
	{Label: "Sarapan 🍳", StartHour: 6, StartMin: 0, EndHour: 9, EndMin: 0},
	{Label: "Makan Siang 🍽️", StartHour: 11, StartMin: 30, EndHour: 13, EndMin: 30},
	{Label: "Makan Malam 🌙", StartHour: 18, StartMin: 0, EndHour: 21, EndMin: 0},
}

// scoredPlace holds a place and its computed score for selection.
type scoredPlace struct {
	place domain.Place
	score float64
	dist  float64
}

// getMealLabel checks if a unix timestamp falls within a meal window.
// Returns the meal label if it does, otherwise empty string.
func getMealLabel(unixTime int64) string {
	t := time.Unix(unixTime, 0)
	minuteOfDay := t.Hour()*60 + t.Minute()

	for _, mw := range mealWindows {
		startMinute := mw.StartHour*60 + mw.StartMin
		endMinute := mw.EndHour*60 + mw.EndMin
		if minuteOfDay >= startMinute && minuteOfDay < endMinute {
			return mw.Label
		}
	}
	return ""
}

// isRestaurantCategory checks if a place belongs to the restaurant category.
func isRestaurantCategory(p *domain.Place) bool {
	cat := strings.ToLower(p.Category)
	return cat == "restaurants"
}

func getAverageRating(itinerary []domain.ScheduledPlace) float64 {
	if len(itinerary) == 0 {
		return 0
	}
	sum := 0.0
	for _, item := range itinerary {
		sum += item.Place.Rating
	}
	avg := sum / float64(len(itinerary))
	return math.Round(avg*1000) / 1000
}

type CandidateItinerary struct {
	Itinerary     []domain.ScheduledPlace
	TotalDistance float64
	TotalDuration int64
	AvgRating     float64
	Score         float64 // Lower is better
	Geometry      interface{}
}

func (pu *planningUsecase) GenerateRecommendation(c context.Context, req *domain.PlanningRequest) (domain.PlanningResponse, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	allPlaces, err := pu.placeRepository.Fetch(ctx)
	if err != nil {
		return domain.PlanningResponse{}, err
	}

	// Define category groups to handle strict filtering
	// Key: Category, Value: Sub-types belonging to it
	typeGroups := map[string][]string{
		"tourist attractions": {"alam", "budaya", "buatan"},
		"restaurants":         {"halal (tersertifikasi)", "halal (belum tersertifikasi)", "vegetarian", "non-halal"},
	}

	// 1. Filter places by preference
	filteredPlaces := []domain.Place{}
	prefMap := make(map[string]bool)
	for _, p := range req.Preferences {
		prefMap[strings.ToLower(p)] = true
	}

	for _, p := range allPlaces {
		if len(req.Preferences) == 0 {
			filteredPlaces = append(filteredPlaces, p)
			continue
		}

		pType := strings.ToLower(p.Type)
		pCat := strings.ToLower(p.Category)

		match := false

		// Check if user selected this specific type
		if prefMap[pType] {
			match = true
		}

		// Check if user selected this broad category
		if !match && prefMap[pCat] {
			hasSpecificTypeSelection := false
			if subTypes, ok := typeGroups[pCat]; ok {
				for _, st := range subTypes {
					if prefMap[st] {
						hasSpecificTypeSelection = true
						break
					}
				}
			}

			if !hasSpecificTypeSelection {
				match = true
			}
		}

		// Special handling for keywords like "halal"
		if !match && prefMap["halal"] {
			if strings.Contains(pType, "halal") || strings.Contains(pCat, "halal") {
				match = true
			}
		}

		if match {
			filteredPlaces = append(filteredPlaces, p)
		}
	}

	fmt.Printf("Total places: %d, Filtered places: %d\n", len(allPlaces), len(filteredPlaces))

	// 2. Separate into restaurants and attractions, group attractions by sub-type
	restaurants := []domain.Place{}
	attractions := []domain.Place{}
	attractionsByType := make(map[string][]domain.Place) // e.g. "alam" -> [...], "budaya" -> [...]
	for _, p := range filteredPlaces {
		if isRestaurantCategory(&p) {
			restaurants = append(restaurants, p)
		} else {
			attractions = append(attractions, p)
			subType := strings.ToLower(p.Type)
			attractionsByType[subType] = append(attractionsByType[subType], p)
		}
	}

	// Build list of selected attraction sub-types
	requestedSubTypes := []string{}
	knownSubTypes := []string{"alam", "budaya", "buatan"}
	for _, st := range knownSubTypes {
		if prefMap[st] || (prefMap["tourist attractions"] && len(attractionsByType[st]) > 0) {
			requestedSubTypes = append(requestedSubTypes, st)
		}
	}
	if len(requestedSubTypes) == 0 && len(attractions) > 0 {
		for _, st := range knownSubTypes {
			if len(attractionsByType[st]) > 0 {
				requestedSubTypes = append(requestedSubTypes, st)
			}
		}
	}

	fmt.Printf("Restaurants: %d, Attractions: %d, Sub-types: %v\n", len(restaurants), len(attractions), requestedSubTypes)

	// 3. Generate Multiple Itineraries for Ranking (Max 50 variations)
	var generatedItineraries []CandidateItinerary

	// Ensure we generate unique paths
	seenPaths := make(map[string]bool)

	for iter := 0; iter < 50; iter++ {
		var rng *rand.Rand
		if iter > 0 {
			rng = rand.New(rand.NewSource(int64(iter)))
		}

		itinerary := []domain.ScheduledPlace{}
		currentTime := req.StartTime
		currentLat := req.StartLat
		currentLong := req.StartLong
		totalDistance := 0.0

		visited := make(map[string]bool)
		usedMealWindows := make(map[string]bool)

		typeCounts := make(map[string]int)

		for len(itinerary) < req.MaxPlaces && currentTime < req.EndTime {
			mealLabel := getMealLabel(currentTime)
			isMealTime := mealLabel != "" && !usedMealWindows[mealLabel]

			var candidatePool []domain.Place
			var activityLabel string
			var spendDuration int64

			if isMealTime && len(restaurants) > 0 {
				candidatePool = restaurants
				activityLabel = mealLabel
				spendDuration = 2700 // 45 minutes
			} else {
				candidatePool = attractions
				activityLabel = "Wisata 🎡"
				spendDuration = 3600 // 60 minutes
			}

			hasUnvisited := false
			for _, p := range candidatePool {
				if !visited[p.ID] {
					hasUnvisited = true
					break
				}
			}
			if !hasUnvisited {
				if isMealTime {
					candidatePool = attractions
					activityLabel = "Wisata 🎡"
					spendDuration = 3600
				} else {
					candidatePool = restaurants
					if ml := getMealLabel(currentTime); ml != "" {
						activityLabel = ml
					} else {
						activityLabel = "Kuliner 🍴"
					}
					spendDuration = 2700
				}
			}

			candidates := []scoredPlace{}
			for i := range candidatePool {
				p := &candidatePool[i]
				if visited[p.ID] {
					continue
				}

				dist := geoutil.Haversine(currentLat, currentLong, p.Latitude, p.Longitude)
				travelTimeSeconds := int64((dist / 30.0) * 3600)
				arrivalTime := currentTime + travelTimeSeconds

				if req.ReturnToStart {
					distToStart := geoutil.Haversine(p.Latitude, p.Longitude, req.StartLat, req.StartLong)
					travelTimeBackSeconds := int64((distToStart / 30.0) * 3600)
					if arrivalTime+spendDuration+travelTimeBackSeconds > req.EndTime {
						continue
					}
				} else {
					if arrivalTime+spendDuration > req.EndTime {
						continue
					}
				}

				// Hilangkan batas maksimal (cap) agar jarak yang lebih jauh selalu mendapat penalti
				normDist := dist / 100.0 
				normRating := (5.0 - p.Rating) / 5.0

				// Berikan bobot yang lebih besar pada jarak (70%) agar algoritma memprioritaskan rute terdekat
				score := (normDist * 0.7) + (normRating * 0.3)
				
				if !isRestaurantCategory(p) {
					pSubType := strings.ToLower(p.Type)
					// Add penalty based on how many times this type has been visited
					// A penalty of 1.0 is huge (since natural score is 0.0-1.0),
					// forcing the algorithm to balance the selection among selected types.
					score += float64(typeCounts[pSubType]) * 1.0
				}

				candidates = append(candidates, scoredPlace{
					place: *p,
					score: score,
					dist:  dist,
				})
			}

			if len(candidates) == 0 {
				if isMealTime {
					currentTime += 1800
					continue
				}
				break
			}

			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].score < candidates[j].score
			})

			var selected scoredPlace
			if rng != nil && len(candidates) > 1 {
				topN := 5
				if topN > len(candidates) {
					topN = len(candidates)
				}
				idx := rng.Intn(topN)
				selected = candidates[idx]
			} else {
				selected = candidates[0]
			}

			actualDist := geoutil.Haversine(currentLat, currentLong, selected.place.Latitude, selected.place.Longitude)
			travelTimeSeconds := int64((actualDist / 30.0) * 3600)
			arrivalTime := currentTime + travelTimeSeconds
			departureTime := arrivalTime + spendDuration

			itinerary = append(itinerary, domain.ScheduledPlace{
				Place:         selected.place,
				ArrivalTime:   arrivalTime,
				DepartureTime: departureTime,
				Distance:      actualDist,
				ActivityLabel: activityLabel,
			})

			visited[selected.place.ID] = true
			if isMealTime {
				usedMealWindows[mealLabel] = true
			} else {
				pSubType := strings.ToLower(selected.place.Type)
				typeCounts[pSubType]++
			}
			currentTime = departureTime
			currentLat = selected.place.Latitude
			currentLong = selected.place.Longitude
			totalDistance += actualDist
		}

		if req.ReturnToStart && len(itinerary) > 0 {
			distToStart := geoutil.Haversine(currentLat, currentLong, req.StartLat, req.StartLong)
			travelTimeBackSeconds := int64((distToStart / 30.0) * 3600)
			arrivalTime := currentTime + travelTimeBackSeconds
			
			returnLocationName := "Kembali ke Lokasi Awal"
			if req.StartLocationName != "" {
				returnLocationName = "Kembali ke " + req.StartLocationName
			}

			itinerary = append(itinerary, domain.ScheduledPlace{
				Place: domain.Place{
					ID:        "return-location",
					Name:      returnLocationName,
					Latitude:  req.StartLat,
					Longitude: req.StartLong,
					Type:      "End Point",
					Category:  "End Point",
				},
				ArrivalTime:   arrivalTime,
				DepartureTime: arrivalTime,
				Distance:      distToStart,
				ActivityLabel: "Selesai Perjalanan 🏁",
			})
			currentTime = arrivalTime
			totalDistance += distToStart
		}

		if len(itinerary) > 0 {
			// Create path hash to ensure uniqueness
			pathID := ""
			for _, item := range itinerary {
				pathID += item.Place.ID + "-"
			}
			if !seenPaths[pathID] {
				seenPaths[pathID] = true
				avgRating := getAverageRating(itinerary)
				totalDuration := currentTime - req.StartTime

				generatedItineraries = append(generatedItineraries, CandidateItinerary{
					Itinerary:     itinerary,
					TotalDistance: totalDistance,
					TotalDuration: totalDuration,
					AvgRating:     avgRating,
				})
			}
		}
	}

	if len(generatedItineraries) == 0 {
		return domain.PlanningResponse{Itinerary: []domain.ScheduledPlace{}}, nil
	}

	// WaitGroup and Mutex for concurrent Mapbox requests
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Ambil data Mapbox semua itinerary
	for i := range generatedItineraries {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			it := generatedItineraries[idx]

			coords := []geoutil.Coordinate{
				{Lat: req.StartLat, Long: req.StartLong},
			}
			for _, item := range it.Itinerary {
				coords = append(coords, geoutil.Coordinate{Lat: item.Place.Latitude, Long: item.Place.Longitude})
			}

			mapboxRoute, err := geoutil.GetMapboxRoute(coords)
			if err == nil && mapboxRoute != nil && len(mapboxRoute.Legs) == len(it.Itinerary) {
				runningTime := req.StartTime
				totalDistMapbox := 0.0

				for j := 0; j < len(it.Itinerary); j++ {
					leg := mapboxRoute.Legs[j]
					distKm := leg.Distance / 1000.0
					travelTimeSeconds := int64(leg.Duration)

					arrivalTime := runningTime + travelTimeSeconds
					spendDuration := it.Itinerary[j].DepartureTime - it.Itinerary[j].ArrivalTime
					departureTime := arrivalTime + spendDuration

					it.Itinerary[j].Distance = distKm
					it.Itinerary[j].ArrivalTime = arrivalTime
					it.Itinerary[j].DepartureTime = departureTime

					runningTime = departureTime
					totalDistMapbox += distKm
				}

				it.TotalDistance = totalDistMapbox
				it.TotalDuration = runningTime - req.StartTime
				it.Geometry = mapboxRoute.Geometry
			}

			mu.Lock()
			generatedItineraries[idx] = it
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	// Cari Min dan Max
	minDist := math.MaxFloat64
	maxDist := 0.0
	minDur := int64(math.MaxInt64)
	maxDur := int64(0)
	minRating := math.MaxFloat64
	maxRating := 0.0

	for _, it := range generatedItineraries {
		if it.TotalDistance < minDist {
			minDist = it.TotalDistance
		}
		if it.TotalDistance > maxDist {
			maxDist = it.TotalDistance
		}
		if it.TotalDuration < minDur {
			minDur = it.TotalDuration
		}
		if it.TotalDuration > maxDur {
			maxDur = it.TotalDuration
		}
		if it.AvgRating < minRating {
			minRating = it.AvgRating
		}
		if it.AvgRating > maxRating {
			maxRating = it.AvgRating
		}
	}

	epsilon := 0.0001
	distRange := (maxDist - minDist) + epsilon
	durRange := float64(maxDur - minDur) + epsilon
	ratingRange := (maxRating - minRating) + epsilon

	// Normalisasi dan hitung score (Higher is better)
	for i, it := range generatedItineraries {
		distScore := (maxDist - it.TotalDistance) / distRange
		durScore := float64(maxDur - it.TotalDuration) / durRange
		ratingScore := (it.AvgRating - minRating) / ratingRange

		combinedScore := (distScore * 0.5) + (durScore * 0.3) + (ratingScore * 0.2)
		generatedItineraries[i].Score = combinedScore
	}

	// Sort generated itineraries by Score (best first, highest is best)
	sort.Slice(generatedItineraries, func(i, j int) bool {
		return generatedItineraries[i].Score > generatedItineraries[j].Score
	})

	// Select the requested variation (based on Seed index)
	selectedIndex := req.Seed
	if selectedIndex >= len(generatedItineraries) {
		selectedIndex = len(generatedItineraries) - 1 // Wrap to the worst if out of bounds
	}

	finalItinerary := generatedItineraries[selectedIndex]

	startLocationName := "Lokasi Awal"
	if req.StartLocationName != "" {
		startLocationName = req.StartLocationName
	}

	startLocationPlace := domain.ScheduledPlace{
		Place: domain.Place{
			ID:        "start-location",
			Name:      startLocationName,
			Latitude:  req.StartLat,
			Longitude: req.StartLong,
			Type:      "Starting Point",
			Category:  "Starting Point",
		},
		ArrivalTime:   0,
		DepartureTime: req.StartTime,
		Distance:      0,
		ActivityLabel: "Mulai Perjalanan 🚗",
	}

	finalItineraryList := append([]domain.ScheduledPlace{startLocationPlace}, finalItinerary.Itinerary...)

	return domain.PlanningResponse{
		Itinerary:     finalItineraryList,
		TotalDistance: finalItinerary.TotalDistance,
		TotalDuration: finalItinerary.TotalDuration,
		AverageRating: finalItinerary.AvgRating,
		Geometry:      finalItinerary.Geometry,
	}, nil
}