package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"business-report-system/internal/model"
	"business-report-system/internal/repository"
)

type POSService struct {
	posRepo      *repository.POSCredentialRepo
	reportRepo   *repository.ReportRepo
	businessRepo *repository.BusinessRepo
}

func NewPOSService(pr *repository.POSCredentialRepo, rr *repository.ReportRepo, br *repository.BusinessRepo) *POSService {
	return &POSService{posRepo: pr, reportRepo: rr, businessRepo: br}
}

func (s *POSService) CreateCredential(businessID uint) (*model.POSCredential, error) {
	apiKey, _ := generateRandomHex(32)
	secret, _ := generateRandomHex(64)

	cred := &model.POSCredential{
		BusinessID: businessID,
		APIKey:     apiKey,
		Secret:     secret,
		Status:     "active",
	}
	if err := s.posRepo.Create(cred); err != nil {
		return nil, err
	}
	return cred, nil
}

func (s *POSService) UpdateCredential(id uint, status string) error {
	return s.posRepo.Update(&model.POSCredential{ID: id, Status: status})
}

func (s *POSService) ListCredentials(page, pageSize int) ([]model.POSCredential, int64, error) {
	return s.posRepo.List(page, pageSize)
}

func (s *POSService) SubmitReport(businessID uint, req *POSReportReq) (*model.Report, error) {
	month := time.Now().Format("2006-01")

	reported, err := s.reportRepo.BusinessMonthReported(businessID, month)
	if err != nil {
		return nil, err
	}
	if reported {
		return nil, errors.New("本月已上报，不可重复提交")
	}

	total := req.RestaurantRevenue + req.RetailRevenue + req.AccommodationRevenue +
		req.TobaccoAlcoholRevenue + req.OtherRevenue

	report := &model.Report{
		BusinessID:            businessID,
		ReportMonth:           month,
		RestaurantRevenue:     req.RestaurantRevenue,
		RetailRevenue:         req.RetailRevenue,
		AccommodationRevenue:  req.AccommodationRevenue,
		TobaccoAlcoholRevenue: req.TobaccoAlcoholRevenue,
		OtherRevenue:          req.OtherRevenue,
		TotalRevenue:          total,
		SourceType:            "POS",
		Status:                "submitted",
		SubmittedAt:           time.Now(),
	}

	if err := s.reportRepo.Create(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *POSService) GetReport(id uint) (*model.Report, error) {
	return s.reportRepo.FindByID(id)
}

func (s *POSService) BusinessStatus(businessID uint) (map[string]any, error) {
	month := time.Now().Format("2006-01")
	reported, _ := s.reportRepo.BusinessMonthReported(businessID, month)
	report, _ := s.reportRepo.FindByBusinessAndMonth(businessID, month)

	result := map[string]any{
		"month":    month,
		"reported": reported,
	}
	if reported && report != nil {
		result["report"] = report
	}
	return result, nil
}

func (s *POSService) LogCall(credentialID uint, method, path, body, ip string, statusCode int, durationMs int64) {
	s.posRepo.CreateCallLog(&model.APICallLog{
		CredentialID: credentialID,
		Method:       method,
		Path:         path,
		RequestBody:  body,
		ResponseCode: statusCode,
		IPAddress:    ip,
		DurationMs:   durationMs,
	})
}

func (s *POSService) CallLogs(credentialID uint, page, pageSize int) ([]model.APICallLog, int64, error) {
	return s.posRepo.ListCallLogs(page, pageSize, credentialID)
}

type POSReportReq struct {
	RestaurantRevenue     float64 `json:"restaurant_revenue"`
	RetailRevenue         float64 `json:"retail_revenue"`
	AccommodationRevenue  float64 `json:"accommodation_revenue"`
	TobaccoAlcoholRevenue float64 `json:"tobacco_alcohol_revenue"`
	OtherRevenue          float64 `json:"other_revenue"`
}

func generateRandomHex(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func ComputeSignature(apiKey, timestamp, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(apiKey + ":" + timestamp))
	return hex.EncodeToString(mac.Sum(nil))
}
