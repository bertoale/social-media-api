package user

import (
	"fmt"
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/response"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
	config  *config.Config
}

// helper function to parse uint param from URL
func ParseUserID(c *gin.Context) (uint, error) {
	idParam := c.Param("user_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
func ParseUsername(c *gin.Context) (string, error) {
	usernameParam := c.Param("username")
	if usernameParam == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	return usernameParam, nil
}
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	uid, ok := userID.(uint)
	return uid, ok
}

// Register godoc
// @Summary Register user
// @Description Register
// @Tags User
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "Registration data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/register [post]
func (ctrl *Controller) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := ctrl.service.Register(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "registration successful", user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags User
// @Accept json
// @Produce json
// @Param data body LoginRequest true "Login credentials"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/login [post]
func (ctrl *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	token, user, err := ctrl.service.Login(&req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.SetCookie(
		"token",
		token,
		int((7 * 24 * time.Hour).Seconds()),
		"/",
		"",
		ctrl.config.NodeEnv == "production",
		true,
	)
	response.Success(c, http.StatusOK, "login successful", gin.H{
		"user": user,
		// "token": token,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user information by user ID
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/users/{user_id} [get]
func (ctrl *Controller) GetUserByID(c *gin.Context) {
	userID, err := ParseUserID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID")
		return
	}
	user, err := ctrl.service.GetUserByID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "user fetched successfully", user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's profile information
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param username formData string false "Username"
// @Param bio formData string false "User bio"
// @Param avatar formData file false "Avatar image"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/users/me [put]
func (ctrl *Controller) UpdateProfile(c *gin.Context) {
	authUserID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	oldUser, err := ctrl.service.GetUserByID(authUserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "user not found")
		return
	}

	oldAvatar := oldUser.Avatar

	if uploadedFile, exists := c.Get("uploadedFile"); exists {
		fileStr := uploadedFile.(string)
		if fileStr != "" {
			req.Avatar = &fileStr
		}
	}

	user, err := ctrl.service.UpdateProfile(authUserID, &req)
	if err != nil {
		if err.Error() == "username already in use" || err.Error() == "email already in use" {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Delete old avatar file if a new one was uploaded
	if req.Avatar != nil && oldAvatar != "" && oldAvatar != *req.Avatar {
		_ = os.Remove("." + oldAvatar)
	}

	response.Success(c, http.StatusOK, "profile updated successfully", user)
}

// GetUserByUsername godoc
// @Summary Get user by username
// @Description Retrieve user information by username
// @Tags User
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/users/username/{username} [get]
func (ctrl *Controller) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := ctrl.service.GetUserByUsername(username)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "user fetched successfully", user)
}

// GetCurrentUser godoc
// @Summary Get current authenticated user
// @Description Retrieve information of the currently authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/users/me [get]
func (ctrl *Controller) GetCurrentUser(c *gin.Context) {
	authUserID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := ctrl.service.GetCurrentUserDetail(authUserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "user not found")
		return
	}

	response.Success(c, http.StatusOK, "current user fetched successfully", user)
}

// GetCurrentUser godoc
// @Summary Get current authenticated user
// @Description Retrieve information of the currently authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/users/me [get]
func (ctrl *Controller) GetExploreUsers(c *gin.Context) {
	authUserID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		response.Error(c, http.StatusBadRequest, "invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		response.Error(c, http.StatusBadRequest, "invalid offset parameter")
		return
	}

	users, err := ctrl.service.GetExploreUsers(authUserID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "explore users fetched successfully", users)
}

// GetUserDetailByID godoc
// @Summary Get user detail by ID
// @Description Retrieve detailed user information by user ID
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/users/{user_id}/detail [get]
func (ctrl *Controller) GetUserDetailByUsername(c *gin.Context) {
	username, err := ParseUsername(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	authUserID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := ctrl.service.GetUserDetailByUsername(username, authUserID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user detail fetched successfully", user)
}

func (ctrl *Controller) SearchUser(c *gin.Context) {
	keyword := c.Query("q")
	authUserID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	users, err := ctrl.service.SearchUser(keyword, authUserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "users fetched successfully", users)
}

func (ctrl *Controller) GetUserFollowers(c *gin.Context) {
	authUserID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		response.Error(c, http.StatusBadRequest, "invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		response.Error(c, http.StatusBadRequest, "invalid offset parameter")
		return
	}

	users, err := ctrl.service.GetUserFollowers(authUserID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user followers fetched successfully", users)
}

func (ctrl *Controller) GetUserFollowings(c *gin.Context) {
	authUserID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		response.Error(c, http.StatusBadRequest, "invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		response.Error(c, http.StatusBadRequest, "invalid offset parameter")
		return
	}

	users, err := ctrl.service.GetUserFollowings(authUserID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user followings fetched successfully", users)
}

func (ctrl *Controller) Logout(c *gin.Context) {
	err := ctrl.service.Logout()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to logout")
		return
	}

	// Clear the token cookie
	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		"",
		ctrl.config.NodeEnv == "production",
		true,
	)

	response.Success(c, http.StatusOK, "logout successful", nil)
}

func NewController(s Service, cfg *config.Config) *Controller {
	return &Controller{
		service: s,
		config:  cfg,
	}
}
