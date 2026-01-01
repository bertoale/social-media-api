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
	protectedAPI.PUT("/users/me", middlewares.UploadAvatar(), ctrl.UpdateProfile)
	protectedAPI.GET("/users/me", ctrl.GetCurrentUser)
	protectedAPI.GET("/users/username/:username", ctrl.GetUserDetailByUsername)

	protectedAPI.GET("/users/explore", ctrl.GetExploreUsers)
	protectedAPI.GET("/users/search", ctrl.SearchUser)
	protectedAPI.GET("/users/followers", ctrl.GetUserFollowers)
	protectedAPI.GET("/users/followings", ctrl.GetUserFollowings)
}
