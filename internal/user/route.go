package user

import (
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config) {
	publicAPI := r.Group("/api")
	publicAPI.POST("/register", ctrl.Register)
	publicAPI.POST("/login", ctrl.Login)

	protectedAPI := r.Group("/api")
	protectedAPI.Use(middlewares.Authenticate(cfg))
	protectedAPI.GET("/users/:user_id", ctrl.GetUserByID)
	protectedAPI.PUT("/users/me", middlewares.UploadAvatar(), ctrl.UpdateProfile)
}
