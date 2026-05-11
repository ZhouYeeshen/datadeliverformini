package handler

import (
	"net/http"
	"strconv"

	"business-report-system/internal/service"

	"github.com/gin-gonic/gin"
)

type BusinessHandler struct {
	svc *service.BusinessService
}

func NewBusinessHandler(svc *service.BusinessService) *BusinessHandler {
	return &BusinessHandler{svc: svc}
}

func (h *BusinessHandler) GetProfile(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未绑定企业"})
		return
	}

	profile, err := h.svc.GetProfile(businessID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "企业不存在"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *BusinessHandler) UpdateProfile(c *gin.Context) {
	businessID := c.GetUint("business_id")
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.svc.UpdateBusiness(businessID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *BusinessHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	list, total, err := h.svc.List(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total, "page": page, "page_size": pageSize})
}

func (h *BusinessHandler) GetDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := h.svc.GetDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "企业不存在"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (h *BusinessHandler) Create(c *gin.Context) {
	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if biz, err := h.svc.Create(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"business": biz})
	}
}

func (h *BusinessHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.svc.UpdateBusiness(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
