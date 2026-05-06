package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BusinessContext(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType := c.GetString("user_type")
		userID := c.GetUint("user_id")

		if userType == "wechat" {
			var result struct {
				BusinessID uint
			}
			if err := db.Table("wechat_accounts").
				Select("business_id").
				Where("id = ? AND status = 'active'", userID).
				First(&result).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "请先绑定企业"})
				c.Abort()
				return
			}
			c.Set("business_id", result.BusinessID)
		}

		if userType == "admin" {
			c.Set("business_id", uint(0))
		}

		c.Next()
	}
}
