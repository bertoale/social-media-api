package report

import (
	"go-sosmed/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func ParseReportID(c *gin.Context) (uint, error) {
	idParam := c.Param("report_id")
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

func ParsePostID(c *gin.Context) (uint, error) {
	idParam := c.Param("post_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// CreateReport godoc
// @Summary Create a report
// @Description Report a post for inappropriate content
// @Tags Report
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param data body ReportRequest true "Report data"
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/posts/{post_id}/reports [post]
func (ctrl *Controller) CreateReport(c *gin.Context) {
	var req ReportRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	postID, err := ParsePostID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post ID: "+err.Error())
		return
	}
	resp, err := ctrl.service.CreateReport(userID, postID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"report created successfully",
		resp,
	)
}

// GetReportByID godoc
// @Summary Get report by ID
// @Description Retrieve a specific report by ID (Admin only)
// @Tags Report
// @Accept json
// @Produce json
// @Param report_id path int true "Report ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/reports/{report_id} [get]
func (ctrl *Controller) GetReportByID(c *gin.Context) {
	id, err := ParseReportID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid report ID: "+err.Error())
		return
	}
	resp, err := ctrl.service.GetReportByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "report retrieved successfully", resp)
}

// GetAllReports godoc
// @Summary Get all reports
// @Description Retrieve all reports (Admin only)
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/reports [get]
func (ctrl *Controller) GetAllReports(c *gin.Context) {
	resps, err := ctrl.service.GetAllReports()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "reports retrieved successfully", resps)
}

// UpdateReportStatus godoc
// @Summary Update report status
// @Description Update the status of a report (Admin only)
// @Tags Report
// @Accept json
// @Produce json
// @Param report_id path int true "Report ID"
// @Param data body UpdateReportRequest true "Status update data"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /api/reports/{report_id}/status [put]
func (ctrl *Controller) UpdateReportStatus(c *gin.Context) {
	id, err := ParseReportID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid report ID: "+err.Error())
		return
	}
	var req UpdateReportRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	resp, err := ctrl.service.UpdateReportStatus(id, req.Status)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "report status updated successfully", resp)
}

func NewController(service Service) *Controller {
	return &Controller{service}
}
