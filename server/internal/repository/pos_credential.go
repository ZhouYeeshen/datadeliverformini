package repository

import (
	"business-report-system/internal/model"

	"gorm.io/gorm"
)

type POSCredentialRepo struct{ db *gorm.DB }

func NewPOSCredentialRepo(db *gorm.DB) *POSCredentialRepo { return &POSCredentialRepo{db: db} }

func (r *POSCredentialRepo) Create(c *model.POSCredential) error {
	return r.db.Create(c).Error
}

func (r *POSCredentialRepo) FindByBusinessID(businessID uint) (*model.POSCredential, error) {
	var c model.POSCredential
	err := r.db.Where("business_id = ?", businessID).First(&c).Error
	return &c, err
}

func (r *POSCredentialRepo) FindByAPIKey(apiKey string) (*model.POSCredential, error) {
	var c model.POSCredential
	err := r.db.Where("api_key = ?", apiKey).First(&c).Error
	return &c, err
}

func (r *POSCredentialRepo) Update(c *model.POSCredential) error {
	return r.db.Save(c).Error
}

func (r *POSCredentialRepo) List(page, pageSize int) ([]model.POSCredential, int64, error) {
	var list []model.POSCredential
	var total int64

	q := r.db.Model(&model.POSCredential{})
	q.Count(&total)

	err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error
	return list, total, err
}

func (r *POSCredentialRepo) CreateCallLog(log *model.APICallLog) error {
	return r.db.Create(log).Error
}

func (r *POSCredentialRepo) ListCallLogs(page, pageSize int, credentialID uint) ([]model.APICallLog, int64, error) {
	var list []model.APICallLog
	var total int64

	q := r.db.Model(&model.APICallLog{})
	if credentialID > 0 {
		q = q.Where("credential_id = ?", credentialID)
	}
	q.Count(&total)

	err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error
	return list, total, err
}
