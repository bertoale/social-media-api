package report

func ToReportResponse(report *Report) *ReportResponse {
	return &ReportResponse{
		ID:        report.ID,
		UserID:    report.UserID,
		BlogID:    report.BlogID,
		Reason:    report.Reason,
		Status:    report.Status,
		CreatedAt: report.CreatedAt,
		UpdatedAt: report.UpdatedAt,
	}
}
