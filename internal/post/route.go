package post

import (
	"go-sosmed/internal/session"
	cloudinaryPkg "go-sosmed/pkg/cloudinary"
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupPostRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config, sessionService session.Service, cloudinaryService *cloudinaryPkg.Service) {
	postGroup := r.Group("/api/posts")
	{
		postGroup.GET("", ctrl.GetAllUnarchived)
		postGroup.GET("/:post_id", middlewares.Authenticate(cfg, sessionService), ctrl.GetDetailByID)
		postGroup.POST("", middlewares.Authenticate(cfg, sessionService), middlewares.UploadPostImage(cloudinaryService), ctrl.Create)
		postGroup.PUT("/:post_id", middlewares.Authenticate(cfg, sessionService), middlewares.UploadPostImage(cloudinaryService), ctrl.Update)
		postGroup.DELETE("/:post_id", middlewares.Authenticate(cfg, sessionService), ctrl.Delete)
		postGroup.GET("/author/:author_id", ctrl.GetPostsByAuthor)
		postGroup.GET("/author/me", middlewares.Authenticate(cfg, sessionService), ctrl.GetAllByCurrentUser)
		postGroup.PATCH("/:post_id/archive", middlewares.Authenticate(cfg, sessionService), ctrl.Archive)
		postGroup.PATCH("/:post_id/unarchive", middlewares.Authenticate(cfg, sessionService), ctrl.Unarchive)
		postGroup.GET("/following", middlewares.Authenticate(cfg, sessionService), ctrl.GetPostsByFollowing)
		postGroup.GET("/liked/me", middlewares.Authenticate(cfg, sessionService), ctrl.GetLikedPosts)
	}
}
