package utils

import (
	"testing"
	"time"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Email válido",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "Email inválido - sem @",
			email:    "testexample.com",
			expected: false,
		},
		{
			name:     "Email inválido - sem domínio",
			email:    "test@",
			expected: false,
		},
		{
			name:     "Email inválido - sem extensão",
			email:    "test@example",
			expected: false,
		},
		{
			name:     "Email inválido - extensão curta",
			email:    "test@example.c",
			expected: false,
		},
		{
			name:     "Email vazio",
			email:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%s) = %v; esperava %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected bool
	}{
		{
			name:     "Username válido",
			username: "test_user",
			expected: true,
		},
		{
			name:     "Username muito curto",
			username: "ab",
			expected: false,
		},
		{
			name:     "Username muito longo",
			username: "a_very_long_username_that_exceeds_twenty_characters",
			expected: false,
		},
		{
			name:     "Username com caracteres inválidos",
			username: "test-user",
			expected: false,
		},
		{
			name:     "Username vazio",
			username: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUsername(tt.username)
			if result != tt.expected {
				t.Errorf("IsValidUsername(%s) = %v; esperava %v", tt.username, result, tt.expected)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Senha válida",
			password: "password123",
			expected: true,
		},
		{
			name:     "Senha muito curta",
			password: "pass",
			expected: false,
		},
		{
			name:     "Senha com 8 caracteres",
			password: "password",
			expected: true,
		},
		{
			name:     "Senha vazia",
			password: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPassword(tt.password)
			if result != tt.expected {
				t.Errorf("IsValidPassword(%s) = %v; esperava %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestIsFutureDate(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{
			name:     "Data no futuro",
			date:     time.Now().Add(24 * time.Hour),
			expected: true,
		},
		{
			name:     "Data no passado",
			date:     time.Now().Add(-24 * time.Hour),
			expected: false,
		},
		{
			name:     "Data atual",
			date:     time.Now(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFutureDate(tt.date)
			if result != tt.expected {
				t.Errorf("IsFutureDate(%v) = %v; esperava %v", tt.date, result, tt.expected)
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		list     []string
		expected bool
	}{
		{
			name:     "String encontrada",
			str:      "apple",
			list:     []string{"apple", "banana", "orange"},
			expected: true,
		},
		{
			name:     "String não encontrada",
			str:      "grape",
			list:     []string{"apple", "banana", "orange"},
			expected: false,
		},
		{
			name:     "String encontrada (case insensitive)",
			str:      "APPLE",
			list:     []string{"apple", "banana", "orange"},
			expected: true,
		},
		{
			name:     "Lista vazia",
			str:      "apple",
			list:     []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringInSlice(tt.str, tt.list)
			if result != tt.expected {
				t.Errorf("StringInSlice(%s, %v) = %v; esperava %v", tt.str, tt.list, result, tt.expected)
			}
		})
	}
}