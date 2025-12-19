package report

func ToReportResponse(report *Report) *ReportResponse {
	return &ReportResponse{
		ID:        report.ID,
		UserID:    report.UserID,
		PostID:    report.PostID,
		Reason:    report.Reason,
		Status:    report.Status,
		CreatedAt: report.CreatedAt,
		UpdatedAt: report.UpdatedAt,
	}
}
