package handler

import (
	"net/http"

	"business-report-system/internal/config"
	"business-report-system/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
	cfg *config.Config
}

func NewAuthHandler(svc *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{svc: svc, cfg: cfg}
}

type wechatLoginReq struct {
	Code string `json:"code" binding:"required"`
}

func (h *AuthHandler) WeChatLogin(c *gin.Context) {
	var req wechatLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	resp, err := h.svc.WeChatLogin(req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) WeChatBind(c *gin.Context) {
	openID := c.GetString("openid")
	if openID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证信息"})
		return
	}

	var req service.BindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		gin.DefaultWriter.Write([]byte("[DEBUG bind] openid=" + openID + " bindErr=" + err.Error() + "\n"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 手机号快速关联：只传 phone(+real_name)，无需企业信息
	if req.Phone != "" && req.BusinessName == "" {
		resp, err := h.svc.AutoBindByPhone(openID, req.Phone, req.RealName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if resp != nil {
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	// 完整注册需要校验企业必填字段
	if req.BusinessName == "" || req.LegalPerson == "" || req.IndustryType == "" || req.RealName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写完整的注册信息"})
		return
	}

	resp, err := h.svc.Bind(openID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

type adminLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req adminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	token, err := h.svc.AdminLogin(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
