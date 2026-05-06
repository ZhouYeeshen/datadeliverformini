package repository

import (
	"business-report-system/internal/model"
	"time"

	"gorm.io/gorm"
)

type ReportRepo struct{ db *gorm.DB }

func NewReportRepo(db *gorm.DB) *ReportRepo { return &ReportRepo{db: db} }

func (r *ReportRepo) Create(report *model.Report) error {
	return r.db.Create(report).Error
}

func (r *ReportRepo) FindByBusinessAndMonth(businessID uint, month string) (*model.Report, error) {
	var report model.Report
	err := r.db.Where("business_id = ? AND report_month = ?", businessID, month).First(&report).Error
	return &report, err
}

func (r *ReportRepo) FindByID(id uint) (*model.Report, error) {
	var report model.Report
	err := r.db.First(&report, id).Error
	return &report, err
}

func (r *ReportRepo) FindByBusinessID(businessID uint, months []string) ([]model.Report, error) {
	var list []model.Report
	q := r.db.Where("business_id = ?", businessID)
	if len(months) > 0 {
		q = q.Where("report_month IN ?", months)
	}
	err := q.Order("report_month DESC").Find(&list).Error
	return list, err
}

func (r *ReportRepo) Update(report *model.Report) error {
	return r.db.Save(report).Error
}

func (r *ReportRepo) AdminList(page, pageSize int, businessName, month, status string) ([]model.Report, int64, error) {
	var list []model.Report
	var total int64

	q := r.db.Model(&model.Report{}).Joins("JOIN businesses ON businesses.id = reports.business_id")
	if businessName != "" {
		q = q.Where("LOWER(businesses.name) LIKE LOWER(?)", "%"+businessName+"%")
	}
	if month != "" {
		q = q.Where("reports.report_month = ?", month)
	}
	if status != "" {
		q = q.Where("reports.status = ?", status)
	}
	q.Count(&total)

	err := q.Select("reports.*, businesses.name as business_name").
		Order("reports.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error
	return list, total, err
}

func (r *ReportRepo) MonthlyTotal(month string) (float64, error) {
	var result struct{ Total float64 }
	err := r.db.Model(&model.Report{}).
		Select("COALESCE(SUM(total_revenue), 0) as total").
		Where("report_month = ? AND status = 'submitted'", month).
		Scan(&result).Error
	return result.Total, err
}

func (r *ReportRepo) BusinessMonthReported(businessID uint, month string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Report{}).
		Where("business_id = ? AND report_month = ? AND status = 'submitted'", businessID, month).
		Count(&count).Error
	return count > 0, err
}

func (r *ReportRepo) GetMonthlyReports(month string) ([]model.Report, error) {
	var list []model.Report
	err := r.db.Where("report_month = ? AND status = 'submitted'", month).Find(&list).Error
	return list, err
}

func (r *ReportRepo) GetPreviousMonthsAvg(businessID uint, month string, count int) (float64, error) {
	parsedMonth, err := time.Parse("2006-01", month)
	if err != nil {
		return 0, err
	}

	var months []string
	for i := 1; i <= count; i++ {
		prev := parsedMonth.AddDate(0, -i, 0)
		months = append(months, prev.Format("2006-01"))
	}

	var result struct{ Avg float64 }
	err = r.db.Model(&model.Report{}).
		Select("COALESCE(AVG(total_revenue), 0) as avg").
		Where("business_id = ? AND report_month IN ? AND status = 'submitted'", businessID, months).
		Scan(&result).Error
	return result.Avg, err
}

func (r *ReportRepo) IndustryTotalsByMonth(month string) (map[string]float64, error) {
	type row struct {
		IndustryType string
		Total        float64
	}
	var rows []row
	err := r.db.Model(&model.Report{}).
		Select("businesses.industry_type, COALESCE(SUM(reports.total_revenue), 0) as total").
		Joins("JOIN businesses ON businesses.id = reports.business_id").
		Where("reports.report_month = ? AND reports.status = 'submitted'", month).
		Group("businesses.industry_type").
		Find(&rows).Error

	result := make(map[string]float64)
	for _, r := range rows {
		result[r.IndustryType] = r.Total
	}
	return result, err
}
