package report

import "fmt"

type Service interface {
	CreateReport(userID uint, postID uint, req *ReportRequest) (*ReportResponse, error)
	GetReportByID(id uint) (*ReportResponse, error)
	GetAllReports() ([]*ReportResponse, error)
	UpdateReportStatus(id uint, status string) (*ReportResponse, error)
}

type service struct {
	repo Repository
}

// CreateReport implements Service.
func (s *service) CreateReport(userID uint, postID uint, req *ReportRequest) (*ReportResponse, error) {
	report := &Report{
		UserID: userID,
		PostID: postID,
		Reason: req.Reason,
		Status: StatusPending,
	}

	if err := s.repo.Create(report); err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return ToReportResponse(report), nil
}

// GetAllReports implements Service.
func (s *service) GetAllReports() ([]*ReportResponse, error) {
	report, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	var responses []*ReportResponse
	for _, r := range report {
		responses = append(responses, ToReportResponse(r))
	}
	return responses, nil
}

// GetReportByID implements Service.
func (s *service) GetReportByID(id uint) (*ReportResponse, error) {
	report, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("report not found: %w", err)
	}
	return ToReportResponse(report), nil
}

// UpdateReportStatus implements Service.
func (s *service) UpdateReportStatus(id uint, status string) (*ReportResponse, error) {
	report, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("report not found: %w", err)
	}
	if status != StatusReviewed && status != StatusResolved && status != StatusRejected {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	report.Status = status

	if err := s.repo.Update(report); err != nil {
		return nil, fmt.Errorf("failed to update report: %w", err)
	}
	return ToReportResponse(report), nil
}

func NewService(repo Repository) Service {
	return &service{repo}
}
