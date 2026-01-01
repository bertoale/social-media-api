package follow

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupFollowRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	api := r.Group("/api/follow")
	api.Use(middlewares.Authenticate(cfg))

	api.POST("/:following_id", ctrl.FollowUser)
	api.DELETE("/:following_id", ctrl.UnfollowUser)

}
