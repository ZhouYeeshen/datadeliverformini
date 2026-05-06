package service

import (
	"errors"

	"business-report-system/internal/config"
	"business-report-system/internal/middleware"
	"business-report-system/internal/model"
	"business-report-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg         *config.Config
	wechatRepo  *repository.WeChatAccountRepo
	businessRepo *repository.BusinessRepo
	adminRepo   *repository.AdminRepo
}

func NewAuthService(cfg *config.Config, wr *repository.WeChatAccountRepo, br *repository.BusinessRepo, ar *repository.AdminRepo) *AuthService {
	return &AuthService{cfg: cfg, wechatRepo: wr, businessRepo: br, adminRepo: ar}
}

type WeChatLoginResp struct {
	Token        string `json:"token"`
	IsBound      bool   `json:"is_bound"`
	BusinessID   uint   `json:"business_id,omitempty"`
	BusinessName string `json:"business_name,omitempty"`
	RealName     string `json:"real_name,omitempty"`
}

const DevTestOpenID = "DEV_TEST_USER_DEV_TEST_USER_DEV"

func (s *AuthService) WeChatLogin(openID string) (*WeChatLoginResp, error) {
	// Dev mode: auto-provision test user
	if openID == DevTestOpenID {
		return s.devTestLogin()
	}

	account, err := s.wechatRepo.FindByOpenID(openID)
	resp := &WeChatLoginResp{IsBound: false}

	if err == nil && account != nil {
		resp.IsBound = true
		resp.BusinessID = account.BusinessID
		resp.RealName = account.RealName

		token, err := middleware.GenerateToken(s.cfg.JWT, account.ID, openID, "wechat")
		if err != nil {
			return nil, err
		}
		resp.Token = token

		biz, err := s.businessRepo.FindByID(account.BusinessID)
		if err == nil {
			resp.BusinessName = biz.Name
		}
		return resp, nil
	}

	// Not bound yet — generate a temporary token for binding
	token, err := middleware.GenerateToken(s.cfg.JWT, 0, openID, "wechat")
	if err != nil {
		return nil, err
	}
	resp.Token = token
	return resp, nil
}

func (s *AuthService) devTestLogin() (*WeChatLoginResp, error) {
	// Ensure test account exists
	account, err := s.wechatRepo.FindByOpenID(DevTestOpenID)
	if err != nil || account == nil {
		// Auto-provision test business and account
		biz, _, err := s.businessRepo.FindOrCreateByName("测试企业-沙县小吃", "张三", "restaurant")
		if err != nil {
			return nil, err
		}
		account = &model.WeChatAccount{
			OpenID:     DevTestOpenID,
			BusinessID: biz.ID,
			RealName:   "张测试",
			Phone:      "13800138000",
			Status:     "active",
		}
		if err := s.wechatRepo.Create(account); err != nil {
			// May already exist from race, re-fetch
			account, err = s.wechatRepo.FindByOpenID(DevTestOpenID)
			if err != nil {
				return nil, err
			}
		}
	}

	token, err := middleware.GenerateToken(s.cfg.JWT, account.ID, DevTestOpenID, "wechat")
	if err != nil {
		return nil, err
	}

	biz, err := s.businessRepo.FindByID(account.BusinessID)
	bizName := ""
	if err == nil {
		bizName = biz.Name
	}

	return &WeChatLoginResp{
		Token:        token,
		IsBound:      true,
		BusinessID:   account.BusinessID,
		BusinessName: bizName,
		RealName:     account.RealName,
	}, nil
}

type BindRequest struct {
	OpenID       string `json:"openid" binding:"required"`
	BusinessName string `json:"business_name" binding:"required"`
	LegalPerson  string `json:"legal_person" binding:"required"`
	IndustryType string `json:"industry_type" binding:"required"`
	RealName     string `json:"real_name" binding:"required"`
	Phone        string `json:"phone"`
	LicenseNo    string `json:"license_no"`
}

func (s *AuthService) Bind(req *BindRequest) (*WeChatLoginResp, error) {
	// Check if already bound
	existing, err := s.wechatRepo.FindByOpenID(req.OpenID)
	if err == nil && existing != nil {
		return nil, errors.New("该微信已绑定企业")
	}

	var biz *model.Business

	// 1. Try to find an existing business by phone number
	if req.Phone != "" {
		found, err := s.businessRepo.FindByPhone(req.Phone)
		if err == nil && found != nil {
			biz = found
		}
	}

	// 2. Try to find by matching business name
	if biz == nil {
		found, _, err := s.businessRepo.FindOrCreateByName(req.BusinessName, req.LegalPerson, req.IndustryType)
		if err != nil {
			return nil, err
		}
		biz = found
	}

	// Update contact phone & license on existing business
	if req.Phone != "" {
		biz.ContactPhone = req.Phone
	} else {
		req.Phone = biz.ContactPhone
	}
	if req.LicenseNo != "" {
		biz.LicenseNo = req.LicenseNo
	}
	s.businessRepo.Update(biz)

	account := &model.WeChatAccount{
		OpenID:     req.OpenID,
		BusinessID: biz.ID,
		RealName:   req.RealName,
		Phone:      req.Phone,
		Status:     "active",
	}
	if err := s.wechatRepo.Create(account); err != nil {
		return nil, err
	}

	token, err := middleware.GenerateToken(s.cfg.JWT, account.ID, req.OpenID, "wechat")
	if err != nil {
		return nil, err
	}

	return &WeChatLoginResp{
		Token:        token,
		IsBound:      true,
		BusinessID:   biz.ID,
		BusinessName: biz.Name,
	}, nil
}

// AutoBindByPhone 通过手机号快速关联已有企业，返回nil表示未找到匹配企业
func (s *AuthService) AutoBindByPhone(openID, phone, realName string) (*WeChatLoginResp, error) {
	// Check if already bound
	existing, err := s.wechatRepo.FindByOpenID(openID)
	if err == nil && existing != nil {
		return nil, errors.New("该微信已绑定企业")
	}

	// Look up business by phone
	biz, err := s.businessRepo.FindByPhone(phone)
	if err != nil || biz == nil {
		return nil, nil // 未找到，由上层继续完整绑定流程
	}

	// Auto-bind wechat account to existing business
	account := &model.WeChatAccount{
		OpenID:     openID,
		BusinessID: biz.ID,
		RealName:   realName,
		Phone:      phone,
		Status:     "active",
	}
	if err := s.wechatRepo.Create(account); err != nil {
		return nil, err
	}

	token, err := middleware.GenerateToken(s.cfg.JWT, account.ID, openID, "wechat")
	if err != nil {
		return nil, err
	}

	return &WeChatLoginResp{
		Token:        token,
		IsBound:      true,
		BusinessID:   biz.ID,
		BusinessName: biz.Name,
	}, nil
}

func (s *AuthService) AdminLogin(username, password string) (string, error) {
	admin, err := s.adminRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("账号或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("账号或密码错误")
	}

	return middleware.GenerateToken(s.cfg.JWT, admin.ID, "", "admin")
}
