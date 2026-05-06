package controller

import (
	"net/http"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/gin-gonic/gin"
)

type PlanningController struct {
	PlanningUsecase domain.PlanningUsecase
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
