package services

import (
	"context"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetTodaySummary() (*models.Report, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.repo.GetSalesSummary(ctx, start, end)
}

func (s *ReportService) GetSummaryByDateRange(
	start, end time.Time,
) (*models.Report, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.repo.GetSalesSummary(ctx, start, end)
}
