package handler

import (
	"net/http"
	"strconv"

	"business-report-system/internal/service"

	"github.com/gin-gonic/gin"
)

type POSHandler struct {
	svc *service.POSService
}

func NewPOSHandler(svc *service.POSService) *POSHandler {
	return &POSHandler{svc: svc}
}

func (h *POSHandler) SubmitReport(c *gin.Context) {
	businessID := c.GetUint("business_id")

	var req service.POSReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	report, err := h.svc.SubmitReport(businessID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *POSHandler) GetReport(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	report, err := h.svc.GetReport(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *POSHandler) BusinessStatus(c *gin.Context) {
	businessID := c.GetUint("business_id")
	statusData, err := h.svc.BusinessStatus(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, statusData)
}

func (h *POSHandler) CreateCredential(c *gin.Context) {
	var req struct {
		BusinessID uint `json:"business_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	cred, err := h.svc.CreateCredential(req.BusinessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cred)
}

func (h *POSHandler) UpdateCredential(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.svc.UpdateCredential(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *POSHandler) ListCredentials(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.svc.ListCredentials(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}

func (h *POSHandler) CallLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	credID, _ := strconv.Atoi(c.DefaultQuery("credential_id", "0"))

	list, total, err := h.svc.CallLogs(uint(credID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}
