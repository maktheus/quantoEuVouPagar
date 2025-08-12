package utils

import (
	"testing"
)

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{
			name:     "Valor inteiro",
			value:    1234,
			expected: "R$ 1.234,00",
		},
		{
			name:     "Valor com centavos",
			value:    1234.56,
			expected: "R$ 1.234,56",
		},
		{
			name:     "Valor com um dígito decimal",
			value:    1234.5,
			expected: "R$ 1.234,50",
		},
		{
			name:     "Valor menor que 1000",
			value:    123.45,
			expected: "R$ 123,45",
		},
		{
			name:     "Valor menor que 100",
			value:    12.34,
			expected: "R$ 12,34",
		},
		{
			name:     "Valor menor que 10",
			value:    1.23,
			expected: "R$ 1,23",
		},
		{
			name:     "Valor menor que 1",
			value:    0.12,
			expected: "R$ 0,12",
		},
		{
			name:     "Valor zero",
			value:    0,
			expected: "R$ 0,00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCurrency(tt.value)
			if result != tt.expected {
				t.Errorf("FormatCurrency(%f) = %s; esperava %s", tt.value, result, tt.expected)
			}
		})
	}
}