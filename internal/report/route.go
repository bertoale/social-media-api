package report

import (
	"go-sosmed/internal/session"
	"go-sosmed/pkg/config"
	"go-sosmed/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine, ctrl *Controller, cfg *config.Config, sessionService session.Service) {

	api := r.Group("/api")
	api.Use(middlewares.Authenticate(cfg, sessionService))

	api.POST("/posts/:post_id/reports", ctrl.CreateReport)
	api.GET("/reports/:report_id", middlewares.Authorize("admin"), ctrl.GetReportByID)
	api.GET("/reports", middlewares.Authorize("admin"), ctrl.GetAllReports)
	api.PUT("/reports/:report_id/status", middlewares.Authorize("admin"), ctrl.UpdateReportStatus)
}
