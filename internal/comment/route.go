package comment

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCommentRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {

	api := r.Group("/api")
	api.Use(middlewares.Authenticate(cfg))

	// Create & get comments of a blog
	api.POST("/blogs/:blog_id/comments", ctrl.CreateComment)
	api.GET("/blogs/:blog_id/comments", ctrl.GetCommentTree)

	// Replies (level 2 rule)
	api.POST("/blogs/:blog_id/comments/:id/reply", ctrl.ReplyToComment)
	api.GET("/blogs/:blog_id/comments/:id/replies", ctrl.GetReplies)

	// Comment actions
	api.PUT("/comments/:id", ctrl.UpdateComment)
	api.DELETE("/comments/:id", ctrl.DeleteComment)
}
