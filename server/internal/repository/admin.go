package repository

import (
	"business-report-system/internal/model"

	"gorm.io/gorm"
)

type AdminRepo struct{ db *gorm.DB }

func NewAdminRepo(db *gorm.DB) *AdminRepo { return &AdminRepo{db: db} }

func (r *AdminRepo) FindByUsername(username string) (*model.Admin, error) {
	var a model.Admin
	err := r.db.Where("username = ? AND status = 'active'", username).First(&a).Error
	return &a, err
}

func (r *AdminRepo) Create(a *model.Admin) error {
	return r.db.Create(a).Error
}
