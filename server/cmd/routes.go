package main

import (
	"business-report-system/internal/config"
	"business-report-system/internal/handler"
	"business-report-system/internal/middleware"
	"business-report-system/internal/repository"
	"business-report-system/internal/scheduler"
	"business-report-system/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppServices struct {
	AlertSvc *service.AlertService
}

func registerRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config) *AppServices {
	businessRepo := repository.NewBusinessRepo(db)
	wechatRepo := repository.NewWeChatAccountRepo(db)
	reportRepo := repository.NewReportRepo(db)
	alertRepo := repository.NewAlertRepo(db)
	adminRepo := repository.NewAdminRepo(db)
	posRepo := repository.NewPOSCredentialRepo(db)

	authSvc := service.NewAuthService(cfg, wechatRepo, businessRepo, adminRepo)
	businessSvc := service.NewBusinessService(businessRepo, wechatRepo)
	reportSvc := service.NewReportService(reportRepo, businessRepo)
	alertSvc := service.NewAlertService(alertRepo, reportRepo, businessRepo)
	statsSvc := service.NewStatisticsService(reportRepo, businessRepo)
	posSvc := service.NewPOSService(posRepo, reportRepo, businessRepo)
	ocrSvc := service.NewOCRService(cfg.WeChat.AppID, cfg.WeChat.AppSecret,
		cfg.TencentCloud.SecretID, cfg.TencentCloud.SecretKey, cfg.TencentCloud.Region)

	authH := handler.NewAuthHandler(authSvc, cfg)
	businessH := handler.NewBusinessHandler(businessSvc)
	reportH := handler.NewReportHandler(reportSvc)
	alertH := handler.NewAlertHandler(alertSvc)
	statsH := handler.NewStatisticsHandler(statsSvc)
	posH := handler.NewPOSHandler(posSvc)
	ocrH := handler.NewOCRHandler(ocrSvc)

	jwtAuth := middleware.AuthMiddleware(cfg.JWT)
	bizCtx := middleware.BusinessContext(db)
	posAuth := middleware.POSAuthMiddleware(db)

	auth := r.Group("/auth")
	{
		auth.POST("/wechat/login", authH.WeChatLogin)
		auth.POST("/wechat/bind", jwtAuth, authH.WeChatBind)
		auth.POST("/admin/login", authH.AdminLogin)
	}

	wx := r.Group("/wx")
	wx.Use(jwtAuth, bizCtx)
	{
		wx.GET("/profile", businessH.GetProfile)
		wx.PUT("/profile", businessH.UpdateProfile)
		wx.POST("/reports", reportH.Submit)
		wx.GET("/reports/current", reportH.CurrentMonth)
		wx.GET("/reports/history", reportH.History)
		wx.GET("/reports/:id", reportH.GetDetail)
		wx.POST("/upload/photo", reportH.UploadPhoto)
		wx.POST("/ocr/recognize", ocrH.Recognize)
		wx.POST("/ocr/parse", ocrH.ParseText)
	}

	admin := r.Group("/admin")
	admin.Use(jwtAuth, bizCtx)
	{
		admin.GET("/businesses", businessH.List)
		admin.POST("/businesses", businessH.Create)
		admin.GET("/businesses/:id", businessH.GetDetail)
		admin.PUT("/businesses/:id", businessH.AdminUpdate)
		admin.GET("/businesses/:id/ledger", reportH.Ledger)
		admin.GET("/reports", reportH.AdminList)
		admin.GET("/reports/:id", reportH.AdminGetDetail)
		admin.PUT("/reports/:id", reportH.AdminUpdate)
		admin.GET("/alerts", alertH.List)
		admin.GET("/alerts/unread-count", alertH.UnreadCount)
		admin.PUT("/alerts/:id/read", alertH.MarkRead)
		admin.GET("/statistics/overview", statsH.Overview)
		admin.GET("/statistics/trend", statsH.Trend)
		admin.GET("/statistics/industry", statsH.ByIndustry)
		admin.GET("/statistics/monthly", statsH.MonthlySummary)
		admin.POST("/pos/credentials", posH.CreateCredential)
		admin.PUT("/pos/credentials/:id", posH.UpdateCredential)
		admin.GET("/pos/credentials", posH.ListCredentials)
		admin.GET("/pos/logs", posH.CallLogs)
	}

	pos := r.Group("/pos")
	pos.Use(posAuth)
	{
		pos.POST("/report", posH.SubmitReport)
		pos.GET("/report/:id", posH.GetReport)
		pos.GET("/status", posH.BusinessStatus)
	}

	return &AppServices{AlertSvc: alertSvc}
}

func startScheduler(svcs *AppServices) {
	s := scheduler.New(svcs.AlertSvc)
	s.Start()
}
