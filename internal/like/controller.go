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

// helper parse blog ID
func ParseBlogID(c *gin.Context) (uint, error) {
	idParam := c.Param("blog_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
func ParseUserID(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
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

func (ctrl *Controller) LikeBlog(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	blogID, err := ParseBlogID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid blog ID")
		return
	}
	err = ctrl.service.LikeBlog(userID, blogID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "blog liked successfully", nil)
}

func (ctrl *Controller) UnlikeBlog(c *gin.Context) {
	blogID, err := ParseBlogID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid blog ID")
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	err = ctrl.service.UnlikeBlog(userID, blogID)
	if err != nil {
		switch {
		case err.Error() == "not liked yet":
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "blog unliked successfully", nil)
}

func (ctrl *Controller) GetBlogLikedByCurrentUser(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return
	}

	blogs, err := ctrl.service.GetBlogsLikedByUser(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "liked blogs retrieved successfully", blogs)
}

func (ctrl *Controller) GetBlogLikedByUser(c *gin.Context) {
	userID, err := ParseUserID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID")
		return
	}
	blogs, err := ctrl.service.GetBlogsLikedByUser(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "liked blogs retrieved successfully", blogs)
}

func (ctrl *Controller) IsBlogLiked(c *gin.Context) {
	blogID, err := ParseBlogID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid blog ID")
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	liked, err := ctrl.service.IsBlogLiked(userID, blogID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "like status retrieved successfully", gin.H{"liked": liked})
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}
