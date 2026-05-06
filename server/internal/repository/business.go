package repository

import (
	"business-report-system/internal/model"

	"gorm.io/gorm"
)

type BusinessRepo struct{ db *gorm.DB }

func NewBusinessRepo(db *gorm.DB) *BusinessRepo { return &BusinessRepo{db: db} }

func (r *BusinessRepo) Create(b *model.Business) error {
	return r.db.Create(b).Error
}

func (r *BusinessRepo) FindByID(id uint) (*model.Business, error) {
	var b model.Business
	err := r.db.Where("id = ?", id).First(&b).Error
	return &b, err
}

func (r *BusinessRepo) FindByPhone(phone string) (*model.Business, error) {
	var b model.Business
	err := r.db.Where("contact_phone = ? AND status = 'active'", phone).First(&b).Error
	return &b, err
}

func (r *BusinessRepo) FindByName(name string) (*model.Business, error) {
	var b model.Business
	err := r.db.Where("name = ? AND status = 'active'", name).First(&b).Error
	return &b, err
}

func (r *BusinessRepo) FindOrCreateByName(name, legalPerson, industryType string) (*model.Business, bool, error) {
	var b model.Business
	err := r.db.Where("name = ?", name).First(&b).Error
	if err == nil {
		return &b, false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, false, err
	}

	b = model.Business{
		Name:        name,
		LegalPerson: legalPerson,
		IndustryType: industryType,
		Status:      "active",
	}
	if err := r.db.Create(&b).Error; err != nil {
		return nil, false, err
	}
	return &b, true, nil
}

func (r *BusinessRepo) List(page, pageSize int, keyword string) ([]model.Business, int64, error) {
	var list []model.Business
	var total int64

	q := r.db.Model(&model.Business{})
	if keyword != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?) OR LOWER(legal_person) LIKE LOWER(?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	q.Count(&total)

	err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error
	return list, total, err
}

func (r *BusinessRepo) Update(b *model.Business) error {
	return r.db.Save(b).Error
}

func (r *BusinessRepo) AllActive() ([]model.Business, error) {
	var list []model.Business
	err := r.db.Where("status = 'active'").Find(&list).Error
	return list, err
}

func (r *BusinessRepo) Count() (int64, error) {
	var c int64
	err := r.db.Model(&model.Business{}).Where("status = 'active'").Count(&c).Error
	return c, err
}

func (r *BusinessRepo) CountByIndustry() (map[string]int64, error) {
	type row struct {
		IndustryType string
		Count        int64
	}
	var rows []row
	err := r.db.Model(&model.Business{}).
		Select("industry_type, count(*) as count").
		Where("status = 'active'").
		Group("industry_type").
		Find(&rows).Error

	result := make(map[string]int64)
	for _, r := range rows {
		result[r.IndustryType] = r.Count
	}
	return result, err
}
