package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"business-report-system/internal/service"

	"github.com/gin-gonic/gin"
)

type OCRHandler struct {
	svc *service.OCRService
}

func NewOCRHandler(svc *service.OCRService) *OCRHandler {
	return &OCRHandler{svc: svc}
}

type ocrRequest struct {
	ImageBase64 string `json:"image_base64"`
	FilePath    string `json:"file_path"`
}

type ocrParseRequest struct {
	RawText string `json:"raw_text"`
}

type ocrResponse struct {
	RestaurantRevenue     float64 `json:"restaurant_revenue"`
	RetailRevenue         float64 `json:"retail_revenue"`
	AccommodationRevenue  float64 `json:"accommodation_revenue"`
	TobaccoAlcoholRevenue float64 `json:"tobacco_alcohol_revenue"`
	OtherRevenue          float64 `json:"other_revenue"`
	TotalRevenue          float64 `json:"total_revenue"`
	RawText               string  `json:"raw_text"`
	Engine                string  `json:"engine"`
	Confidence            string  `json:"confidence"`
	PhotoURL              string  `json:"photo_url"`
}

func (h *OCRHandler) Recognize(c *gin.Context) {
	var req ocrRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var result *service.OCRResult
	var photoURL string
	var err error

	if req.FilePath != "" {
		result, err = h.svc.ProcessImage(req.FilePath)
		photoURL = req.FilePath
	} else if req.ImageBase64 != "" {
		data, derr := base64.StdEncoding.DecodeString(req.ImageBase64)
		if derr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "图片数据格式错误"})
			return
		}
		os.MkdirAll("uploads", 0755)
		filename := fmt.Sprintf("uploads/ocr_%d.png", time.Now().UnixNano())
		if werr := os.WriteFile(filename, data, 0644); werr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
			return
		}
		photoURL = "/" + filename
		result, err = h.svc.ProcessImageFromBase64(req.ImageBase64)
	} else {
		file, ferr := c.FormFile("photo")
		if ferr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请提供图片"})
			return
		}
		tmpPath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, tmpPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
			return
		}
		photoURL = "/" + tmpPath
		result, err = h.svc.ProcessImage(tmpPath)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OCR识别失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ocrResponse{
		RestaurantRevenue:     result.RestaurantRevenue,
		RetailRevenue:         result.RetailRevenue,
		AccommodationRevenue:  result.AccommodationRevenue,
		TobaccoAlcoholRevenue: result.TobaccoAlcoholRevenue,
		OtherRevenue:          result.OtherRevenue,
		TotalRevenue:          result.TotalRevenue,
		RawText:               result.RawText,
		Engine:                result.Engine,
		Confidence:            result.Confidence,
		PhotoURL:              photoURL,
	})
}

// ParseText parses OCR text (already recognized by mini program plugin) into structured amounts.
func (h *OCRHandler) ParseText(c *gin.Context) {
	var req ocrParseRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RawText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 raw_text"})
		return
	}

	result := h.svc.ParseText(req.RawText)

	c.JSON(http.StatusOK, ocrResponse{
		RestaurantRevenue:     result.RestaurantRevenue,
		RetailRevenue:         result.RetailRevenue,
		AccommodationRevenue:  result.AccommodationRevenue,
		TobaccoAlcoholRevenue: result.TobaccoAlcoholRevenue,
		OtherRevenue:          result.OtherRevenue,
		TotalRevenue:          result.TotalRevenue,
		RawText:               result.RawText,
		Engine:                result.Engine,
		Confidence:            result.Confidence,
		PhotoURL:              "",
	})
}
