package controller

import (
	"net/http"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	ProfileUsecase domain.ProfileUsecase
}

func (pc *ProfileController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	profile, err := pc.ProfileUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (pc *ProfileController) Update(c *gin.Context) {
	userID := c.GetString("x-user-id")

	var request domain.UpdateProfileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user := &domain.User{
		Name: request.Name,
	}

	err = pc.ProfileUsecase.UpdateProfile(c, userID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.UpdateProfileResponse{Message: "Profile updated successfully"})
}
