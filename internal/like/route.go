package like

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupLikeRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	api := r.Group("/api")
	api.Use(middlewares.Authenticate(cfg))
	api.POST("/posts/:post_id/like", ctrl.LikePost)
	api.DELETE("/posts/:post_id/like", ctrl.UnlikePost)
	api.GET("/posts/:post_id/like/status", ctrl.IsPostLiked)
	api.GET("/users/:user_id/likes", ctrl.GetPostLikedByUser)
	api.GET("/users/me/likes", ctrl.GetPostLikedByCurrentUser)
	api.GET("/posts/:post_id/like/count", ctrl.GetPostLikeCount)
}
