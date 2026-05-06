package route

import (
	"database/sql"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/controller"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/repository"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/usecase"
	"github.com/gin-gonic/gin"
)

func NewTripRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, group *gin.RouterGroup) {
	tr := repository.NewTripRepository(db)
	trr := repository.NewTripRouteRepository(db)
	pr := repository.NewPlaceRepository(db)
	tc := &controller.TripController{
		TripUsecase: usecase.NewTripUsecase(tr, trr, pr, timeout),
		Env:         env,
	}
	group.POST("/trips", tc.Create)
	group.GET("/trips", tc.Fetch)
	group.GET("/trips/:id", tc.GetByID)
	group.DELETE("/trips/:id", tc.Delete)
}
