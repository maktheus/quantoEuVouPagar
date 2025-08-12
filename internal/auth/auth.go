package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserRole representa o papel do usuário.
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// UserClaims representa as claims customizadas do JWT.
type UserClaims struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
	jwt.RegisteredClaims
}

// SecretKey é a chave usada para assinar os tokens JWT.
// Em produção, isso deve ser uma variável de ambiente.
var SecretKey []byte

// SetSecretKey define a chave secreta para assinatura de tokens JWT.
func SetSecretKey(key string) {
	SecretKey = []byte(key)
}

// GenerateToken gera um novo token JWT para um usuário.
func GenerateToken(userID int, username string, role UserRole) (string, error) {
	// Define as claims do token
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expira em 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}

	// Cria o token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a chave secreta
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", fmt.Errorf("erro ao assinar o token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida um token JWT e retorna as claims.
func ValidateToken(tokenString string) (*UserClaims, error) {
	// Parseia o token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verifica o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear o token: %w", err)
	}

	// Verifica se o token é válido
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}

// HasPermission verifica se um usuário tem permissão para acessar um recurso.
// Esta é uma implementação básica. Pode ser expandida para verificar permissões mais complexas.
func HasPermission(claims *UserClaims, requiredRole UserRole) bool {
	// Usuário admin tem acesso a tudo
	if claims.Role == RoleAdmin {
		return true
	}

	// Verifica se o papel do usuário é igual ou superior ao papel requerido
	// Neste exemplo simples, apenas "admin" tem mais permissões que "user"
	return claims.Role == requiredRole
}