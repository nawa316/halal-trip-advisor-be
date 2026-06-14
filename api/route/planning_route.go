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

func NewPlanningRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, group *gin.RouterGroup) {
	pr := repository.NewPlaceRepository(db)
	pc := &controller.PlanningController{
		PlanningUsecase: usecase.NewPlanningUsecase(pr, timeout),
		PlaceRepository: pr,
		Env:             env,
	}
	group.POST("/generate", pc.Generate)
	group.GET("/alternatives", pc.GetAlternatives)
}
