// internal/http/handlers/auth/jwt.go
package jwt

import (
	"fmt"
	"time"

	configpkg "goProcessClient/internal/config"
	domain "goProcessClient/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateTokens gera access + refresh tokens para um Vendor
func GenerateTokens(v domain.Vendor) (domain.LoginResponse, error) {
	// 60 minutos para o access token (ajuste se quiser)
	accessExpiresAt := time.Now().Add(60 * time.Minute)

	accessClaims := &domain.UserClaims{
		Cod:     v.Cod,
		Usuario: v.Usuario,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", v.Cod),
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			Issuer:    "goProcessClient",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString(configpkg.JwtSecret)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("erro ao assinar access token: %w", err)
	}

	// Refresh token de 7 dias
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", v.Cod),
		ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
		Issuer:    "goProcessClient",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString(configpkg.JwtSecret)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("erro ao assinar refresh token: %w", err)
	}

	expiresIn := int64(time.Until(accessExpiresAt).Seconds())

	return domain.LoginResponse{
		AccessToken:  signedAccess,
		RefreshToken: signedRefresh,
		ExpiresIn:    expiresIn,
	}, nil
}

// ValidateToken valida um JWT de access token e devolve os claims
func ValidateToken(tokenStr string) (*domain.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("algoritmo de assinatura inválido")
		}
		return configpkg.JwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear token: %w", err)
	}

	claims, ok := token.Claims.(*domain.UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}
