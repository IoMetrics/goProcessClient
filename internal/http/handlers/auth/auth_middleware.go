package auth

import (
	"context"
	domain "goProcessClient/internal/domain"
	jwt "goProcessClient/internal/http/handlers/jwt"
	"net/http"
	"strings"
)

type ctxKey string

const userCtxKey ctxKey = "userClaims"

// AuthMiddleware protege rotas usando o access token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "token ausente")
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "token inválido ou expirado")
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// (Opcional) Helper pra recuperar o UserClaims do contexto
func GetUserClaims(r *http.Request) (*domain.UserClaims, bool) {
	val := r.Context().Value(userCtxKey)
	if val == nil {
		return nil, false
	}
	claims, ok := val.(*domain.UserClaims)
	return claims, ok
}
