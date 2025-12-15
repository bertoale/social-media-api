package user

import (
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
// @Param data body LoginRequest true "Login data"
// @Success 200 {object} response.SuccessResponse
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

// Login handles user login
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
		"user":  user,
		"token": token,
	})
}

// GetUserByID handles fetching user by ID
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

// UpdateProfile handles updating user profile
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
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Delete old avatar file if a new one was uploaded
	if req.Avatar != nil && oldAvatar != "" && oldAvatar != *req.Avatar {
		_ = os.Remove("." + oldAvatar)
	}

	response.Success(c, http.StatusOK, "profile updated successfully", user)
}

func NewController(s Service, cfg *config.Config) *Controller {
	return &Controller{
		service: s,
		config:  cfg,
	}
}
