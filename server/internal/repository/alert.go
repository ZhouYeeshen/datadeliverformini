package repository

import (
	"business-report-system/internal/model"

	"gorm.io/gorm"
)

type AlertRepo struct{ db *gorm.DB }

func NewAlertRepo(db *gorm.DB) *AlertRepo { return &AlertRepo{db: db} }

func (r *AlertRepo) Create(a *model.Alert) error {
	return r.db.Create(a).Error
}

func (r *AlertRepo) List(page, pageSize int, alertType string, isRead *bool) ([]model.Alert, int64, error) {
	var list []model.Alert
	var total int64

	q := r.db.Model(&model.Alert{})
	if alertType != "" {
		q = q.Where("alert_type = ?", alertType)
	}
	if isRead != nil {
		q = q.Where("is_read = ?", *isRead)
	}
	q.Count(&total)

	err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error
	return list, total, err
}

func (r *AlertRepo) UnreadCount() (int64, error) {
	var c int64
	err := r.db.Model(&model.Alert{}).Where("is_read = false").Count(&c).Error
	return c, err
}

func (r *AlertRepo) MarkRead(id uint) error {
	return r.db.Model(&model.Alert{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *AlertRepo) FindByID(id uint) (*model.Alert, error) {
	var a model.Alert
	err := r.db.First(&a, id).Error
	return &a, err
}

func (r *AlertRepo) ExistsByBusinessMonth(businessID uint, alertMonth, alertType string) (bool, error) {
	var c int64
	err := r.db.Model(&model.Alert{}).
		Where("business_id = ? AND alert_month = ? AND alert_type = ?", businessID, alertMonth, alertType).
		Count(&c).Error
	return c > 0, err
}
