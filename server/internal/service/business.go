package service

import (
	"business-report-system/internal/model"
	"business-report-system/internal/repository"
)

type BusinessService struct {
	bizRepo    *repository.BusinessRepo
	wxRepo     *repository.WeChatAccountRepo
}

func NewBusinessService(br *repository.BusinessRepo, wr *repository.WeChatAccountRepo) *BusinessService {
	return &BusinessService{bizRepo: br, wxRepo: wr}
}

type BusinessProfile struct {
	Business    model.Business        `json:"business"`
	Operators   []model.WeChatAccount `json:"operators"`
	ReportedThisMonth bool            `json:"reported_this_month"`
}

func (s *BusinessService) GetProfile(businessID uint) (*BusinessProfile, error) {
	biz, err := s.bizRepo.FindByID(businessID)
	if err != nil {
		return nil, err
	}

	ops, _ := s.wxRepo.FindByBusinessID(businessID)

	return &BusinessProfile{
		Business:  *biz,
		Operators: ops,
	}, nil
}

func (s *BusinessService) UpdateBusiness(businessID uint, updates map[string]any) error {
	biz, err := s.bizRepo.FindByID(businessID)
	if err != nil {
		return err
	}

	if v, ok := updates["name"]; ok {
		biz.Name = v.(string)
	}
	if v, ok := updates["legal_person"]; ok {
		biz.LegalPerson = v.(string)
	}
	if v, ok := updates["industry_type"]; ok {
		biz.IndustryType = v.(string)
	}
	if v, ok := updates["license_no"]; ok {
		biz.LicenseNo = v.(string)
	}
	if v, ok := updates["address"]; ok {
		biz.Address = v.(string)
	}
	if v, ok := updates["contact_phone"]; ok {
		biz.ContactPhone = v.(string)
	}
	if v, ok := updates["status"]; ok {
		biz.Status = v.(string)
	}
	return s.bizRepo.Update(biz)
}

func (s *BusinessService) List(keyword string, page, pageSize int) ([]model.Business, int64, error) {
	return s.bizRepo.List(page, pageSize, keyword)
}

func (s *BusinessService) GetDetail(id uint) (*BusinessProfile, error) {
	return s.GetProfile(id)
}
