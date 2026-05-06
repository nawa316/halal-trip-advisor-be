package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/internal/idutil"
	"github.com/gin-gonic/gin"
)

type TripController struct {
	TripUsecase domain.TripUsecase
	Env         *bootstrap.Env
}

func (tc *TripController) Create(c *gin.Context) {
	var request domain.SaveTripRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID := c.GetString("x-user-id")
	tripID := idutil.NewID()

	routes := make([]domain.TripRoute, len(request.PlaceIDs))
	for i, placeID := range request.PlaceIDs {
		routes[i] = domain.TripRoute{
			ID:         idutil.NewID(),
			TripID:     tripID,
			PlaceID:    placeID,
			OrderIndex: int64(i),
		}
	}

	trip := domain.Trip{
		ID:            tripID,
		UserID:        userID,
		Name:          request.Name,
		StartTime:     request.StartTime,
		EndTime:       request.EndTime,
		TotalDistance: request.TotalDistance,
		Routes:        routes,
	}

	err = tc.TripUsecase.Create(c, &trip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Trip saved successfully",
	})
}

func (tc *TripController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	trips, err := tc.TripUsecase.FetchByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, trips)
}

func (tc *TripController) GetByID(c *gin.Context) {
	tripID := c.Param("id")
	fmt.Printf("Fetching trip with ID: %s\n", tripID)

	trip, err := tc.TripUsecase.GetByID(c, tripID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Trip not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, trip)
}

func (tc *TripController) Delete(c *gin.Context) {
	tripID := c.Param("id")
	fmt.Printf("Deleting trip with ID: %s\n", tripID)

	err := tc.TripUsecase.Delete(c, tripID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Trip deleted successfully",
	})
}
