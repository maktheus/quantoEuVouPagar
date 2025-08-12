package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/auth"
	"go.uber.org/zap"
)

// ContextKey é um tipo para chaves de contexto.
type ContextKey string

const (
	// UserClaimsKey é a chave para armazenar as claims do usuário no contexto.
	UserClaimsKey ContextKey = "user_claims"
)

// AuthMiddleware é um middleware que verifica o token JWT e adiciona as claims do usuário ao contexto.
func AuthMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtém o header de autorização
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("Header de autorização não fornecido", zap.String("path", r.URL.Path))
				http.Error(w, "Authorization header não fornecido", http.StatusUnauthorized)
				return
			}

			// Verifica se o header começa com "Bearer "
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				log.Warn("Header de autorização com prefixo inválido", zap.String("path", r.URL.Path))
				http.Error(w, "Authorization header com formato inválido", http.StatusUnauthorized)
				return
			}

			// Extrai o token
			tokenString := authHeader[len(bearerPrefix):]

			// Valida o token
			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				log.Warn("Token inválido", zap.Error(err), zap.String("path", r.URL.Path))
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			// Adiciona as claims do usuário ao contexto
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaimsFromContext extrai as claims do usuário do contexto.
func GetClaimsFromContext(ctx context.Context) (*auth.UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*auth.UserClaims)
	return claims, ok
}