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
// FollowUser godoc
// @Summary Follow a user
// @Description Follow another user
// @Tags Follow
// @Accept json
// @Produce json
// @Param following_id path int true "User ID to follow"
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/follow/{following_id} [post]
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

// UnfollowUser godoc
// @Summary Unfollow a user
// @Description Unfollow a user you're currently following
// @Tags Follow
// @Accept json
// @Produce json
// @Param following_id path int true "User ID to unfollow"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/follow/{following_id} [delete]
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

// GetFollowers godoc
// @Summary Get followers
// @Description Get list of users following the authenticated user
// @Tags Follow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/follow/me/followers [get]
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

// GetFollowing godoc
// @Summary Get following
// @Description Get list of users the authenticated user is following
// @Tags Follow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/follow/me/following [get]
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
