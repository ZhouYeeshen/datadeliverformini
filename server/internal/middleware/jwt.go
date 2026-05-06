package middleware

import (
	"time"

	"business-report-system/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(cfg config.JWTConfig, userID uint, openID, userType string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		OpenID:   openID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}
