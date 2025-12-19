package comment

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupCommentRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	api := r.Group("/api")
	api.Use(middlewares.Authenticate(cfg))

	// Create & get comments of a post
	api.POST("/posts/:post_id/comments", ctrl.CreateComment)
	api.GET("/posts/:post_id/comments", ctrl.GetCommentTree)

	// Replies (level 2 rule)
	api.POST("/posts/:post_id/comments/:comment_id/reply", ctrl.ReplyToComment)
	api.GET("/posts/:post_id/comments/:comment_id/replies", ctrl.GetReplies)

	// Comment actions
	api.PUT("/comments/:comment_id", ctrl.UpdateComment)
	api.DELETE("/comments/:comment_id", ctrl.DeleteComment)
	api.GET("/posts/:post_id/comments/count", ctrl.GetCommentCount)
}
