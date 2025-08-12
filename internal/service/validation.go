package service

import (
	"fmt"
	"time"

	"github.com/matheus-uchoa/quanto-eu-vou-pagar/internal/model"
)

// ValidarSubconta verifica se os dados da subconta são válidos.
func ValidarSubconta(subconta model.Subconta) error {
	if subconta.Descricao == "" {
		return fmt.Errorf("descrição da subconta é obrigatória")
	}

	if subconta.Valor <= 0 {
		return fmt.Errorf("valor da subconta deve ser maior que zero")
	}

	if subconta.TaxaJuros < 0 {
		return fmt.Errorf("taxa de juros não pode ser negativa")
	}

	if subconta.Parcelas <= 0 {
		return fmt.Errorf("número de parcelas deve ser maior que zero")
	}

	if subconta.DataInicio.IsZero() {
		return fmt.Errorf("data de início é obrigatória")
	}

	// Verifica se a data de início não é no futuro
	// Isso pode ser ajustado conforme a regra de negócio
	// if subconta.DataInicio.After(time.Now()) {
	//	 return fmt.Errorf("data de início não pode ser no futuro")
	// }

	if subconta.Sistema != "price" && subconta.Sistema != "sac" {
		return fmt.Errorf("sistema de amortização '%s' não suportado. Use 'price' ou 'sac'", subconta.Sistema)
	}

	if subconta.Tipo == "" {
		return fmt.Errorf("tipo da subconta é obrigatório")
	}

	return nil
}