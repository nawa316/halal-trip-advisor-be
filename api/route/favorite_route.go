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

func NewFavoriteRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, group *gin.RouterGroup) {
	fr := repository.NewFavoriteRepository(db)
	fc := &controller.FavoriteController{
		FavoriteUsecase: usecase.NewFavoriteUsecase(fr, timeout),
		Env:             env,
	}
	group.POST("/favorites", fc.Create)
	group.GET("/favorites", fc.Fetch)
	group.DELETE("/favorites/:place_id", fc.Delete)
}
