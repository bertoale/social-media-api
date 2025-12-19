package like

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupLikeRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	api := r.Group("/api")
	api.Use(middlewares.Authenticate(cfg))
	// Like / Unlike post
	api.POST("/posts/:post_id/like", ctrl.LikePost)
	api.DELETE("/posts/:post_id/like", ctrl.UnlikePost)
	// Check like status
	api.GET("/posts/:post_id/like/status", ctrl.IsPostLiked)
	// Get all likes by specific user
	api.GET("/users/:user_id/likes", ctrl.GetPostLikedByUser)
	// Get all likes by current user
	api.GET("/users/me/likes", ctrl.GetPostLikedByCurrentUser)
}
