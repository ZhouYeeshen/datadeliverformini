package service

import (
	"time"

	"business-report-system/internal/repository"
)

type StatisticsService struct {
	reportRepo   *repository.ReportRepo
	businessRepo *repository.BusinessRepo
}

func NewStatisticsService(rr *repository.ReportRepo, br *repository.BusinessRepo) *StatisticsService {
	return &StatisticsService{reportRepo: rr, businessRepo: br}
}

type OverviewData struct {
	TotalBusinesses    int64              `json:"total_businesses"`
	ByIndustry         map[string]int64   `json:"by_industry"`
	CurrentMonthTotal  float64            `json:"current_month_total"`
	CurrentMonthReports int64             `json:"current_month_reports"`
	IndustryTotals     map[string]float64 `json:"industry_totals"`
}

func (s *StatisticsService) Overview() (*OverviewData, error) {
	month := time.Now().Format("2006-01")

	totalBiz, _ := s.businessRepo.Count()
	byIndustry, _ := s.businessRepo.CountByIndustry()
	monthTotal, _ := s.reportRepo.MonthlyTotal(month)
	industryTotals, _ := s.reportRepo.IndustryTotalsByMonth(month)

	reports, _ := s.reportRepo.GetMonthlyReports(month)

	return &OverviewData{
		TotalBusinesses:     totalBiz,
		ByIndustry:          byIndustry,
		CurrentMonthTotal:   monthTotal,
		CurrentMonthReports: int64(len(reports)),
		IndustryTotals:      industryTotals,
	}, nil
}

type TrendData struct {
	Month   string  `json:"month"`
	Total   float64 `json:"total"`
	Count   int64   `json:"count"`
}

func (s *StatisticsService) Trend(months int) ([]TrendData, error) {
	var result []TrendData

	now := time.Now()
	for i := months - 1; i >= 0; i-- {
		m := now.AddDate(0, -i, 0).Format("2006-01")
		total, _ := s.reportRepo.MonthlyTotal(m)
		reports, _ := s.reportRepo.GetMonthlyReports(m)

		result = append(result, TrendData{
			Month: m,
			Total: total,
			Count: int64(len(reports)),
		})
	}

	return result, nil
}

func (s *StatisticsService) ByIndustry(month string) (map[string]float64, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	return s.reportRepo.IndustryTotalsByMonth(month)
}

func (s *StatisticsService) MonthlySummary(yearMonth string) (map[string]any, error) {
	month := yearMonth
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	total, _ := s.reportRepo.MonthlyTotal(month)
	reports, _ := s.reportRepo.GetMonthlyReports(month)
	industryTotals, _ := s.reportRepo.IndustryTotalsByMonth(month)

	return map[string]any{
		"month":           month,
		"total_revenue":   total,
		"report_count":    len(reports),
		"industry_totals": industryTotals,
	}, nil
}
