package blog

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupBlogRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	blogGroup := r.Group("/api/blogs")
	{
		blogGroup.GET("", ctrl.GetAll)
		blogGroup.GET("/:blog_id", ctrl.GetByID)
		blogGroup.POST("", middlewares.Authenticate(cfg), middlewares.UploadBlogImage(), ctrl.Create)
		blogGroup.PUT("/:blog_id", middlewares.Authenticate(cfg), middlewares.UploadBlogImage(), ctrl.Update)
		blogGroup.DELETE("/:blog_id", middlewares.Authenticate(cfg), ctrl.Delete)
		blogGroup.GET("/author/me", middlewares.Authenticate(cfg), ctrl.GetBlogsByAuthor)
	}
}
