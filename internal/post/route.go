package post

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupPostRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	postGroup := r.Group("/api/posts")
	{
		postGroup.GET("", ctrl.GetAllUnarchived)
		postGroup.GET("/:post_id", middlewares.Authenticate(cfg), ctrl.GetDetailByID)
		postGroup.POST("", middlewares.Authenticate(cfg), middlewares.UploadPostImage(), ctrl.Create)
		postGroup.PUT("/:post_id", middlewares.Authenticate(cfg), middlewares.UploadPostImage(), ctrl.Update)
		postGroup.DELETE("/:post_id", middlewares.Authenticate(cfg), ctrl.Delete)
		postGroup.GET("/author/:author_id", ctrl.GetPostsByAuthor)
		postGroup.GET("/author/me", middlewares.Authenticate(cfg), ctrl.GetAllByCurrentUser)
		postGroup.PATCH("/:post_id/archive", middlewares.Authenticate(cfg), ctrl.Archive)
		postGroup.PATCH("/:post_id/unarchive", middlewares.Authenticate(cfg), ctrl.Unarchive)
		postGroup.GET("/following", middlewares.Authenticate(cfg), ctrl.GetPostsByFollowing)
		postGroup.GET("/liked/me", middlewares.Authenticate(cfg), ctrl.GetLikedPosts)
	}
}
