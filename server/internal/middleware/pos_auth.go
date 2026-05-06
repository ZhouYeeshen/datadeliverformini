package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func POSAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		timestamp := c.GetHeader("X-Timestamp")
		signature := c.GetHeader("X-Signature")

		if apiKey == "" || timestamp == "" || signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少接口认证参数"})
			c.Abort()
			return
		}

		var cred struct {
			ID         uint
			BusinessID uint
			Secret     string
			Status     string
		}
		if err := db.Table("pos_credentials").
			Select("id, business_id, secret, status").
			Where("api_key = ?", apiKey).
			First(&cred).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的API Key"})
			c.Abort()
			return
		}

		if cred.Status != "active" {
			c.JSON(http.StatusForbidden, gin.H{"error": "API凭证已被禁用"})
			c.Abort()
			return
		}

		expectedSig := computeHMAC(apiKey, timestamp, cred.Secret)
		if signature != expectedSig {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
			c.Abort()
			return
		}

		c.Set("credential_id", cred.ID)
		c.Set("business_id", cred.BusinessID)
		c.Set("user_type", "pos")
		c.Next()

		db.Table("pos_credentials").
			Where("id = ?", cred.ID).
			Update("last_called_at", time.Now())
	}
}

func computeHMAC(apiKey, timestamp, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(apiKey + ":" + timestamp))
	return hex.EncodeToString(mac.Sum(nil))
}
