package utils

import (
	"regexp"
	"strings"
	"time"
)

// IsValidEmail verifica se um email é válido.
func IsValidEmail(email string) bool {
	// Expressão regular simples para validar email
	// Esta é uma validação básica e pode ser aprimorada
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidUsername verifica se um nome de usuário é válido.
// Deve conter apenas letras, números e underscores, com comprimento entre 3 e 20 caracteres.
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// IsValidPassword verifica se uma senha é válida.
// Deve ter pelo menos 8 caracteres.
func IsValidPassword(password string) bool {
	return len(password) >= 8
}

// IsFutureDate verifica se uma data é no futuro.
func IsFutureDate(date time.Time) bool {
	return date.After(time.Now())
}

// StringInSlice verifica se uma string está em um slice de strings.
func StringInSlice(str string, list []string) bool {
	for _, item := range list {
		if strings.EqualFold(item, str) {
			return true
		}
	}
	return false
}