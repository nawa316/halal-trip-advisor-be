package controller

import (
	"net/http"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/internal/idutil"
	"github.com/gin-gonic/gin"
)

type FavoriteController struct {
	FavoriteUsecase domain.FavoriteUsecase
	Env             *bootstrap.Env
}

func (fc *FavoriteController) Create(c *gin.Context) {
	var request struct {
		PlaceID string `json:"place_id" binding:"required"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID := c.GetString("x-user-id")

	favorite := domain.Favorite{
		ID:      idutil.NewID(),
		UserID:  userID,
		PlaceID: request.PlaceID,
	}

	err = fc.FavoriteUsecase.Create(c, &favorite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Place added to favorites",
	})
}

func (fc *FavoriteController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	favorites, err := fc.FavoriteUsecase.FetchByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func (fc *FavoriteController) Delete(c *gin.Context) {
	placeID := c.Param("place_id")
	userID := c.GetString("x-user-id")

	err := fc.FavoriteUsecase.Delete(c, userID, placeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Place removed from favorites",
	})
}
