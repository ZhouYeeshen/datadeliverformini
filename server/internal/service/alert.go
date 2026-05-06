package service

import (
	"fmt"
	"time"

	"business-report-system/internal/model"
	"business-report-system/internal/repository"
)

type AlertService struct {
	alertRepo    *repository.AlertRepo
	reportRepo   *repository.ReportRepo
	businessRepo *repository.BusinessRepo
}

func NewAlertService(ar *repository.AlertRepo, rr *repository.ReportRepo, br *repository.BusinessRepo) *AlertService {
	return &AlertService{alertRepo: ar, reportRepo: rr, businessRepo: br}
}

func (s *AlertService) List(alertType string, isRead *bool, page, pageSize int) ([]model.Alert, int64, error) {
	return s.alertRepo.List(page, pageSize, alertType, isRead)
}

func (s *AlertService) UnreadCount() (int64, error) {
	return s.alertRepo.UnreadCount()
}

func (s *AlertService) MarkRead(id uint) error {
	return s.alertRepo.MarkRead(id)
}

func (s *AlertService) ScanAlerts() error {
	month := time.Now().Format("2006-01")
	businesses, err := s.businessRepo.AllActive()
	if err != nil {
		return err
	}

	for _, biz := range businesses {
		s.checkMissingReport(biz, month)
		s.checkAbnormalChange(biz, month)
	}

	return nil
}

func (s *AlertService) checkMissingReport(biz model.Business, month string) {
	reported, _ := s.reportRepo.BusinessMonthReported(biz.ID, month)
	if reported {
		return
	}

	// Only alert in last 5 days of month
	today := time.Now()
	lastDay := time.Date(today.Year(), today.Month()+1, 0, 0, 0, 0, 0, time.Local)
	if today.Before(lastDay.AddDate(0, 0, -4)) {
		return
	}

	exists, _ := s.alertRepo.ExistsByBusinessMonth(biz.ID, month, "MISSING_REPORT")
	if exists {
		return
	}

	s.alertRepo.Create(&model.Alert{
		BusinessID:   biz.ID,
		BusinessName: biz.Name,
		AlertType:    "MISSING_REPORT",
		AlertMonth:   month,
	})
}

func (s *AlertService) checkAbnormalChange(biz model.Business, month string) {
	report, err := s.reportRepo.FindByBusinessAndMonth(biz.ID, month)
	if err != nil || report.Status != "submitted" {
		return
	}

	avg, err := s.reportRepo.GetPreviousMonthsAvg(biz.ID, month, 2)
	if err != nil || avg == 0 {
		return
	}

	change := (report.TotalRevenue - avg) / avg * 100
	if change >= 30 || change <= -30 {
		exists, _ := s.alertRepo.ExistsByBusinessMonth(biz.ID, month, "ABNORMAL_CHANGE")
		if exists {
			return
		}

		s.alertRepo.Create(&model.Alert{
			BusinessID:     biz.ID,
			BusinessName:   biz.Name,
			AlertType:      "ABNORMAL_CHANGE",
			AlertMonth:     month,
			ReferenceValue: avg,
			CurrentValue:   report.TotalRevenue,
			ChangePercent:  fmtPrec(change, 2),
		})
	}
}

func fmtPrec(f float64, prec int) float64 {
	format := fmt.Sprintf("%%.%df", prec)
	s := fmt.Sprintf(format, f)
	var v float64
	fmt.Sscanf(s, "%f", &v)
	return v
}
