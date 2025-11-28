// internal/domain/auth.go
package domain

import "github.com/golang-jwt/jwt/v5"

// AuthRequest representa um JSON simples enviado pelo app.
// Exemplo: { "ws_chave": "..." }
type AuthRequest struct {
	WsChave string `json:"ws_chave"`
}

// LoginRequest é o payload principal de login:
// { "usuario": "...", "senha": "...", "ws_chave": "..." }
type LoginRequest struct {
	Usuario string `json:"usuario"`
	Senha   string `json:"senha"`
	WsChave string `json:"ws_chave"`
}

// LoginResponse é o que será devolvido após login bem-sucedido.
// Pode ser o que o app Android já espera receber.
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int64  `json:"expiresIn,omitempty"`
}

// UserClaims são os dados que vão dentro do JWT (access token).
type UserClaims struct {
	Cod     int    `json:"cod"`
	Usuario string `json:"usuario"`
	jwt.RegisteredClaims
}
