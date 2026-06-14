package controller

import (
	"net/http"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/gin-gonic/gin"
)

type PlanningController struct {
	PlanningUsecase domain.PlanningUsecase
	PlaceRepository domain.PlaceRepository
	Env             *bootstrap.Env
}

func (pc *PlanningController) Generate(c *gin.Context) {
	var request domain.PlanningRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if request.MaxPlaces <= 0 {
		request.MaxPlaces = 5 // Default
	}

	response, err := pc.PlanningUsecase.GenerateRecommendation(c, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (pc *PlanningController) GetAlternatives(c *gin.Context) {
	placeType := c.Query("type")
	category := c.Query("category")

	allPlaces, err := pc.PlaceRepository.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	alternatives := []domain.Place{}
	for _, p := range allPlaces {
		// If no filters provided, show all
		if placeType == "" && category == "" {
			alternatives = append(alternatives, p)
			continue
		}

		// Match if either matches (OR logic)
		match := false
		if placeType != "" && p.Type == placeType {
			match = true
		}
		if category != "" && p.Category == category {
			match = true
		}

		if match {
			alternatives = append(alternatives, p)
		}
	}

	c.JSON(http.StatusOK, alternatives)
}
