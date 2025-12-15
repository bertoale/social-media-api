package like

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupLikeRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	api := r.Group("/api")
	api.Use(middlewares.Authenticate(cfg))
	// Like / Unlike blog
	api.POST("/blogs/:blog_id/like", ctrl.LikeBlog)
	api.DELETE("/blogs/:blog_id/like", ctrl.UnlikeBlog)
	// Check like status
	api.GET("/blogs/:blog_id/like/status", ctrl.IsBlogLiked)
	// Get all likes by specific user
	api.GET("/users/:user_id/likes", ctrl.GetBlogLikedByUser)
	// Get all likes by current user
	api.GET("/users/me/likes", ctrl.GetBlogLikedByCurrentUser)
}
