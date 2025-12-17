package bill

import (
	"context"
	jwt "goProcessClient/internal/http/handlers/jwt"
	"net/http"
	"strings"
)

type ctxKey string

const userCtxKey ctxKey = "userClaims"

func BillMiddleware(next http.Handler) http.Handler {
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
