package repository

import (
	"business-report-system/internal/model"

	"gorm.io/gorm"
)

type WeChatAccountRepo struct{ db *gorm.DB }

func NewWeChatAccountRepo(db *gorm.DB) *WeChatAccountRepo { return &WeChatAccountRepo{db: db} }

func (r *WeChatAccountRepo) FindByOpenID(openID string) (*model.WeChatAccount, error) {
	var a model.WeChatAccount
	err := r.db.Where("open_id = ? AND status = 'active'", openID).First(&a).Error
	return &a, err
}

func (r *WeChatAccountRepo) Create(a *model.WeChatAccount) error {
	return r.db.Create(a).Error
}

func (r *WeChatAccountRepo) FindByBusinessID(businessID uint) ([]model.WeChatAccount, error) {
	var list []model.WeChatAccount
	err := r.db.Where("business_id = ? AND status = 'active'", businessID).Find(&list).Error
	return list, err
}

func (r *WeChatAccountRepo) CountByBusinessID(businessID uint) (int64, error) {
	var c int64
	err := r.db.Model(&model.WeChatAccount{}).
		Where("business_id = ? AND status = 'active'", businessID).
		Count(&c).Error
	return c, err
}
