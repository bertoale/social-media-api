package like

import (
	"go-sosmed/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

// helper parse post ID
func ParsePostID(c *gin.Context) (uint, error) {
	idParam := c.Param("post_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
func ParseUserID(c *gin.Context) (uint, error) {
	idParam := c.Param("user_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// helper get user ID
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	uid, ok := userID.(uint)
	return uid, ok
}

// LikePost godoc
// @Summary Like a post
// @Description Add a like to a post
// @Tags Like
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/like [post]
func (ctrl *Controller) LikePost(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	postID, err := ParsePostID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post ID")
		return
	}
	err = ctrl.service.LikePost(userID, postID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "post liked successfully", nil)
}

// UnlikePost godoc
// @Summary Unlike a post
// @Description Remove like from a post
// @Tags Like
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/like [delete]
func (ctrl *Controller) UnlikePost(c *gin.Context) {
	postID, err := ParsePostID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post ID")
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	err = ctrl.service.UnlikePost(userID, postID)
	if err != nil {
		switch {
		case err.Error() == "not liked yet":
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "post unliked successfully", nil)
}

// GetPostLikedByCurrentUser godoc
// @Summary Get liked posts by current user
// @Description Retrieve all posts liked by the authenticated user
// @Tags Like
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/users/me/likes [get]
func (ctrl *Controller) GetPostLikedByCurrentUser(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return
	}

	posts, err := ctrl.service.GetPostsLikedByUser(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "liked posts retrieved successfully", posts)
}

// GetPostLikedByUser godoc
// @Summary Get liked posts by specific user
// @Description Retrieve all posts liked by a specific user
// @Tags Like
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/users/{user_id}/likes [get]
func (ctrl *Controller) GetPostLikedByUser(c *gin.Context) {
	userID, err := ParseUserID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID")
		return
	}
	posts, err := ctrl.service.GetPostsLikedByUser(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "liked posts retrieved successfully", posts)
}

// IsPostLiked godoc
// @Summary Check if post is liked
// @Description Check if the authenticated user has liked a specific post
// @Tags Like
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/like/status [get]
func (ctrl *Controller) IsPostLiked(c *gin.Context) {
	postID, err := ParsePostID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post ID")
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	liked, err := ctrl.service.IsPostLiked(userID, postID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "like status retrieved successfully", gin.H{"liked": liked})
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}
