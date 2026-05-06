package route

import (
	"database/sql"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/middleware"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *sql.DB, gin *gin.Engine) {
	gin.Use(middleware.CORSMiddleware())

	publicRouter := gin.Group("")
	// All Public APIs
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)
	NewPlanningRouter(env, timeout, db, publicRouter)

	protectedRouter := gin.Group("")
	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
	NewProfileRouter(env, timeout, db, protectedRouter)
	NewTaskRouter(env, timeout, db, protectedRouter)
	NewTripRouter(env, timeout, db, protectedRouter)
	NewFavoriteRouter(env, timeout, db, protectedRouter)
}
