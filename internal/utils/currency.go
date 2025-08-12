package utils

import (
	"strconv"
	"strings"
)

// FormatCurrency formata um valor float64 como moeda brasileira (BRL).
// Ex: 1234.56 -> "R$ 1.234,56"
func FormatCurrency(value float64) string {
	// Converter o valor para string com duas casas decimais
	strValue := strconv.FormatFloat(value, 'f', 2, 64)

	// Separar a parte inteira e a parte decimal
	parts := strings.Split(strValue, ".")
	integerPart := parts[0]
	decimalPart := "00"
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// Adicionar separadores de milhar à parte inteira
	var formattedInteger strings.Builder
	for i, char := range integerPart {
		// Se não for o primeiro caractere e a posição atual for um múltiplo de 3 a partir do final
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			formattedInteger.WriteRune('.')
		}
		formattedInteger.WriteRune(char)
	}

	// Combinar a parte inteira formatada com a parte decimal
	return "R$ " + formattedInteger.String() + "," + decimalPart
}