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

func ParseBlogID(c *gin.Context) (uint, error) {
	idParam := c.Param("blog_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

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
	blogID, err := ParseBlogID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid blog ID: "+err.Error())
		return
	}
	resp, err := ctrl.service.CreateReport(userID, blogID, &req)
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

func (ctrl *Controller) GetAllReports(c *gin.Context) {
	resps, err := ctrl.service.GetAllReports()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "reports retrieved successfully", resps)
}

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
