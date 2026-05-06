package handler

import (
	"net/http"
	"strconv"

	"business-report-system/internal/service"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	svc *service.StatisticsService
}

func NewStatisticsHandler(svc *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{svc: svc}
}

func (h *StatisticsHandler) Overview(c *gin.Context) {
	data, err := h.svc.Overview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *StatisticsHandler) Trend(c *gin.Context) {
	months, _ := strconv.Atoi(c.DefaultQuery("months", "12"))
	if months > 24 {
		months = 24
	}

	data, err := h.svc.Trend(months)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *StatisticsHandler) ByIndustry(c *gin.Context) {
	month := c.Query("month")
	data, err := h.svc.ByIndustry(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *StatisticsHandler) MonthlySummary(c *gin.Context) {
	month := c.Query("month")
	data, err := h.svc.MonthlySummary(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
