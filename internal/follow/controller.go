package follow

import (
	"go-sosmed/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

// ===============================================
// Helper functions
// =============================================
// Helper function to get user ID from context
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	uid, ok := userID.(uint)
	return uid, ok
}

// Helper function to parse following ID from URL parameter
func ParseFollowingID(c *gin.Context) (uint, error) {
	idParam := c.Param("following_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// ===============================================
// Controller methods
// ===============================================
func (ctrl *Controller) FollowUser(c *gin.Context) {
	followerID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	followingID, err := ParseFollowingID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	if followerID == followingID {
		response.Error(c, http.StatusBadRequest, "cannot follow yourself")
		return
	}
	// call service
	err = ctrl.service.FollowUser(followerID, followingID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Successfully followed user", nil)
}

func (ctrl *Controller) UnfollowUser(c *gin.Context) {
	followerID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	followingID, err := ParseFollowingID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID")
		return
	}
	// call service
	err = ctrl.service.UnfollowUser(followerID, followingID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Successfully unfollowed user", nil)
}

func (ctrl *Controller) GetFollowers(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	followerIDs, err := ctrl.service.GetFollowers(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get followers")
		return
	}
	response.Success(c, http.StatusOK, "Followers retrieved successfully", followerIDs)
}

func (ctrl *Controller) GetFollowing(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	followingIDs, err := ctrl.service.GetFollowing(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get following")
		return
	}
	response.Success(c, http.StatusOK, "Following retrieved successfully", followingIDs)
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}
