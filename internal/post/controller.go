package post

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
func ParsePostID(c *gin.Context) (uint, error) {
	idParam := c.Param("post_id")
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

func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("userRole")
	if !exists {
		return "", false
	}
	role, ok := userRole.(string)
	return role, ok
}

func ParseAuthorID(c *gin.Context) (uint, error) {
	idParam := c.Param("author_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// ======================================================
// Controller methods
// ======================================================
// Create godoc
// @Summary Create a new post
// @Description Create a new post  with optional image upload
// @Tags Post
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Post title"
// @Param content formData string true "Post content"
// @Param image formData file false "Post image"
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/posts [post]
func (ctrl *Controller) Create(c *gin.Context) {
	var req PostRequest
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
	response.Success(c, http.StatusCreated, "post created successfully", resp)
}

// GetByID godoc
// @Summary Get post by ID
// @Description Retrieve a post by its ID
// @Tags Post
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/posts/{post_id} [get]
func (ctrl *Controller) GetByID(c *gin.Context) {
	id, err := ParsePostID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post ID")
		return
	}
	post, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "post fetched successfully", post)
}

// Delete godoc
// @Summary Delete a post
// @Description Delete a post (only author or admin can delete)
// @Tags Post
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/posts/{post_id} [delete]
func (ctrl *Controller) Delete(c *gin.Context) {
	id, err := ParsePostID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post ID")
		return
	}

	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userRole, ok := GetUserRoleFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user role not found")
		return
	}

	err = ctrl.service.Delete(uint(id), userID, userRole)
	if err != nil {
		response.Error(c, http.StatusForbidden, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "post deleted successfully", nil)
}


// GetAll godoc
// @Summary Get all unarchived posts
// @Description Retrieve all unarchived posts
// @Tags Post
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/posts [get]
func (ctrl *Controller) GetAllUnarchived(c *gin.Context) {
	posts, err := ctrl.service.GetAllUnarchived()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "posts retrieved successfully", posts)
}

// GetAllByCurrentUser godoc
// @Summary Get all posts by current user
// @Description Retrieve all posts (including archived) created by the authenticated user
// @Tags Post
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/posts/author/me [get]
func (ctrl *Controller) GetAllByCurrentUser(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	posts, err := ctrl.service.GetPostsByCurrentUser(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "posts retrieved successfully", posts)
}

// GetPostsByAuthor godoc
// @Summary Get posts by specific author
// @Description Retrieve all unarchived posts created by a specific author
// @Tags Post
// @Accept json
// @Produce json
// @Param author_id path int true "Author ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/posts/author/{author_id} [get]
func (ctrl *Controller) GetPostsByAuthor(c *gin.Context) {
	authorID, err := ParseAuthorID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid author ID")
		return
	}
	posts, err := ctrl.service.GetPostsByAuthor(authorID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "posts retrieved successfully", posts)
}

// Update godoc
// @Summary Update a post
// @Description Update a post (only author can update)
// @Tags Post
// @Accept multipart/form-data
// @Produce json
// @Param post_id path int true "Post ID"
// @Param title formData string false "Post title"
// @Param content formData string false "Post content"
// @Param archived formData boolean false "Archive status"
// @Param edited formData boolean false "Edited status"
// @Param image formData file false "Post image"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/posts/{post_id} [put]
func (ctrl *Controller) Update(c *gin.Context) {
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

	var req UpdatePostRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	updatedPost, err := ctrl.service.Update(userID, uint(postID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "post updated successfully", updatedPost)
}

// Archive godoc
// @Summary Archive a post
// @Description Archive a post (only author can archive)
// @Tags Post
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/archive [patch]
func (ctrl *Controller) Archive(c *gin.Context) {
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

	err = ctrl.service.Archive(uint(postID), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "post archived successfully", nil)
}

// Unarchive godoc
// @Summary Unarchive a post
// @Description Unarchive a post (only author can unarchive)
// @Tags Post
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/unarchive [patch]
func (ctrl *Controller) Unarchive(c *gin.Context) {
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

	err = ctrl.service.Unarchive(uint(postID), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "post unarchived successfully", nil)
}
// GetPostsByFollowing godoc
// @Summary Get posts by users the current user is following
// @Description Retrieve all unarchived posts created by users that the authenticated user is following
// @Tags Post
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/posts/following [get]
func (ctrl *Controller) GetPostsByFollowing(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	posts, err := ctrl.service.GetPostsByFollowing(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "posts retrieved successfully", posts)
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}
