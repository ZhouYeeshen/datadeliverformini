package service

import (
	"errors"
	"time"

	"business-report-system/internal/model"
	"business-report-system/internal/repository"
)

type ReportService struct {
	reportRepo   *repository.ReportRepo
	businessRepo *repository.BusinessRepo
}

func NewReportService(rr *repository.ReportRepo, br *repository.BusinessRepo) *ReportService {
	return &ReportService{reportRepo: rr, businessRepo: br}
}

type SubmitReportReq struct {
	BusinessID            uint     `json:"business_id"`
	TotalRevenue          float64  `json:"total_revenue"`
	RestaurantRevenue     float64  `json:"restaurant_revenue"`
	RetailRevenue         float64  `json:"retail_revenue"`
	AccommodationRevenue  float64  `json:"accommodation_revenue"`
	TobaccoAlcoholRevenue float64  `json:"tobacco_alcohol_revenue"`
	OtherRevenue          float64  `json:"other_revenue"`
	SourceType            string   `json:"source_type"`
	PhotoURLs             []string `json:"photo_urls"`
	SubmittedBy           uint     `json:"submitted_by"`
}

func (s *ReportService) Submit(req *SubmitReportReq) (*model.Report, error) {
	month := time.Now().Format("2006-01")

	photosStr := ""
	if len(req.PhotoURLs) > 0 {
		photosStr = req.PhotoURLs[0]
		for i := 1; i < len(req.PhotoURLs); i++ {
			photosStr += "," + req.PhotoURLs[i]
		}
	}

	totalRevenue := req.TotalRevenue
	categorySum := req.RestaurantRevenue + req.RetailRevenue + req.AccommodationRevenue + req.TobaccoAlcoholRevenue + req.OtherRevenue
	if categorySum > 0 && categorySum > totalRevenue {
		totalRevenue = categorySum
	}

	report := &model.Report{
		BusinessID:            req.BusinessID,
		ReportMonth:           month,
		RestaurantRevenue:     req.RestaurantRevenue,
		RetailRevenue:         req.RetailRevenue,
		AccommodationRevenue:  req.AccommodationRevenue,
		TobaccoAlcoholRevenue: req.TobaccoAlcoholRevenue,
		OtherRevenue:          req.OtherRevenue,
		TotalRevenue:          totalRevenue,
		SourceType:            req.SourceType,
		PhotoURLs:             photosStr,
		Status:                "submitted",
		SubmittedBy:           req.SubmittedBy,
		SubmittedAt:           time.Now(),
	}

	if err := s.reportRepo.Create(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *ReportService) CurrentMonthStatus(businessID uint) (map[string]any, error) {
	month := time.Now().Format("2006-01")

	report, err := s.reportRepo.FindByBusinessAndMonth(businessID, month)
	reported := err == nil && report.Status == "submitted"

	result := map[string]any{
		"month":    month,
		"reported": reported,
	}
	if reported {
		result["report"] = report
	}
	return result, nil
}

func (s *ReportService) History(businessID uint) ([]model.Report, error) {
	return s.reportRepo.FindByBusinessID(businessID, nil)
}

func (s *ReportService) GetDetail(id uint) (*model.Report, error) {
	return s.reportRepo.FindByID(id)
}

func (s *ReportService) Ledger(businessID uint) ([]model.Report, error) {
	return s.reportRepo.FindByBusinessID(businessID, nil)
}

func (s *ReportService) AdminList(businessName, month, status string, page, pageSize int) ([]model.Report, int64, error) {
	return s.reportRepo.AdminList(page, pageSize, businessName, month, status)
}

func (s *ReportService) AdminGetDetail(id uint) (*model.Report, error) {
	return s.reportRepo.FindByID(id)
}

type UpdateReportReq struct {
	ReportMonth string  `json:"report_month"`
	TotalRevenue float64 `json:"total_revenue"`
	Status      string  `json:"status"`
	PhotoURLs   string  `json:"photo_urls"`
}

func (s *ReportService) AdminUpdate(id uint, req *UpdateReportReq) (*model.Report, error) {
	report, err := s.reportRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("记录不存在")
	}

	if req.ReportMonth != "" {
		report.ReportMonth = req.ReportMonth
	}
	if req.TotalRevenue > 0 {
		report.TotalRevenue = req.TotalRevenue
	}
	if req.Status != "" {
		report.Status = req.Status
	}
	if req.PhotoURLs != "" {
		report.PhotoURLs = req.PhotoURLs
	}

	if err := s.reportRepo.Update(report); err != nil {
		return nil, err
	}
	return report, nil
}
