package blog

import (
	"go-sosmed/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

// ======================================================
// Helper methods
// ======================================================
// helper function to parse uint param from URL
func ParseBlogID(c *gin.Context) (uint, error) {
	idParam := c.Param("blog_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// helper function to get userID from context
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	uid, ok := userID.(uint)
	return uid, ok
}

// ======================================================
// Controller methods
// ======================================================
func (ctrl *Controller) Create(c *gin.Context) {
	var req BlogRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	// Check for uploaded file (if any)
	if uploadedFile, exists := c.Get("uploadedFile"); exists {
		fileStr := uploadedFile.(string)
		if fileStr != "" {
			req.Image = fileStr
		}
	}
	// Check token
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	// Call service
	resp, err := ctrl.service.Create(&req, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	// Success
	response.Success(c, http.StatusCreated, "blog created successfully", resp)
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	id, err := ParseBlogID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid blog ID")
		return
	}
	blog, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "blog fetched successfully", blog)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	id, err := ParseBlogID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid blog ID")
		return
	}
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	err = ctrl.service.Delete(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "blog deleted successfully", nil)
}

func (ctrl *Controller) GetAll(c *gin.Context) {
	blogs, err := ctrl.service.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "blogs retrieved successfully", blogs)
}

func (ctrl *Controller) GetBlogsByAuthor(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	blogs, err := ctrl.service.GetBlogsByAuthor(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "blogs retrieved successfully", blogs)
}

func (ctrl *Controller) Update(c *gin.Context) {
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

	var req UpdateBlogRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	updatedBlog, err := ctrl.service.Update(userID, uint(blogID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "blog updated successfully", updatedBlog)
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}
