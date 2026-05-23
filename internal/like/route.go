package like

import (
	"go-sosmed/internal/session"
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupLikeRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config, sessionService session.Service) {
	api := r.Group("/api")
	api.Use(middlewares.Authenticate(cfg, sessionService))
	api.POST("/posts/:post_id/like", ctrl.LikePost)
	api.DELETE("/posts/:post_id/like", ctrl.UnlikePost)
	api.GET("/posts/:post_id/like/status", ctrl.IsPostLiked)
	api.GET("/users/:user_id/likes", ctrl.GetPostLikedByUser)
	api.GET("/users/me/likes", ctrl.GetPostLikedByCurrentUser)
}
